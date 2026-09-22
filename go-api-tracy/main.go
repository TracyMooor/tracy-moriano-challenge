package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type NodeRequest struct {
	Matrices [][][]float64 `json:"matrices"`
}

// qrDecomposition realiza la factorización QR mediante Gram-Schmidt modificado
func qrDecomposition(a [][]float64) ([][]float64, [][]float64, error) {
	m := len(a)
	if m == 0 {
		return nil, nil, fmt.Errorf("la matriz no puede estar vacía")
	}
	n := len(a[0])
	if m < n {
		return nil, nil, fmt.Errorf("la matriz debe tener al menos tantas filas como columnas (m >= n)")
	}

	q := make([][]float64, m)
	for i := range q {
		q[i] = make([]float64, n)
	}
	r := make([][]float64, n)
	for i := range r {
		r[i] = make([]float64, n)
	}

	// Copiar columnas para ortogonalizar
	v := make([][]float64, n)
	for j := 0; j < n; j++ {
		v[j] = make([]float64, m)
		for i := 0; i < m; i++ {
			v[j][i] = a[i][j]
		}
	}

	for j := 0; j < n; j++ {
		var norm float64
		for i := 0; i < m; i++ {
			norm += v[j][i] * v[j][i]
		}
		norm = math.Sqrt(norm)
		r[j][j] = norm

		if norm < 1e-12 {
			return nil, nil, fmt.Errorf("las columnas de la matriz no son linealmente independientes")
		}

		// Vector ortonormal de Q
		for i := 0; i < m; i++ {
			q[i][j] = v[j][i] / norm
		}

		// Proyección ortogonal sobre los vectores siguientes
		for k := j + 1; k < n; k++ {
			var dot float64
			for i := 0; i < m; i++ {
				dot += q[i][j] * v[k][i]
			}
			r[j][k] = dot
			for i := 0; i < m; i++ {
				v[k][i] -= dot * q[i][j]
			}
		}
	}

	return q, r, nil
}

func main() {
	app := fiber.New()

	nodeURL := os.Getenv("NODE_API_URL")
	if nodeURL == "" {
		nodeURL = "http://localhost:3000"
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "go-api"})
	})

	app.Post("/api/matrix/qr", func(c *fiber.Ctx) error {
		var req QRRequest
		if err := c.BodyParser(&req); err != nil || len(req.Matrix) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Debe enviar un JSON con el formato: { \"matrix\": [[1, 2], [3, 4]] }",
			})
		}

		// 1. Calcular Factorización QR
		q, r, err := qrDecomposition(req.Matrix)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// 2. Enviar las matrices resultantes a Node.js mediante HTTP
		payload := NodeRequest{Matrices: [][][]float64{q, r}}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error serializando matrices"})
		}

		resp, err := http.Post(nodeURL+"/api/stats", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error":   "Fallo de conexión con la API de Node.js",
				"details": err.Error(),
			})
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error leyendo respuesta de Node.js"})
		}

		var stats map[string]interface{}
		if err := json.Unmarshal(body, &stats); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parseando estadísticas de Node.js"})
		}

		// 3. Devolver la respuesta integrada
		return c.JSON(fiber.Map{
			"q_matrix":   q,
			"r_matrix":   r,
			"statistics": stats,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}