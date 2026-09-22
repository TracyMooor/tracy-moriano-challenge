# Interseguro Coding Challenge Tracy Moriano - Solución Técnica

Solución técnica compuesta por dos microservicios comunicados vía HTTP para procesamiento algebraico y estadístico de matrices.

## Arquitectura
- **go-api (Go + Fiber):** Punto de entrada que recibe la matriz rectangular, calcula la factorización QR ($A = Q \cdot R$) mediante el método de Gram-Schmidt y envía las matrices resultantes al servicio de análisis.
- **node-api (Node.js + Express):** Microservicio analítico que recibe las matrices y extrae: valor máximo, valor mínimo, suma total, promedio y validación de matriz diagonal.

## Ejecución Local con Docker Compose
```bash
docker compose up --build

---

Pruebas de Endpoints
Método: POST
Ruta: http://localhost:8080/api/matrix/qr
Headers: Content-Type: application/json

Payload:
{
  "matrix": [
    [1, 2],
    [3, 4],
    [5, 6]
  ]
}