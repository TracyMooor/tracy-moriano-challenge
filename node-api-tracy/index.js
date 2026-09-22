const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

// Verifica si una matriz cuadrada es diagonal
function isDiagonal(matrix, tol = 1e-7) {
  if (!Array.isArray(matrix) || matrix.length === 0) return false;
  const rows = matrix.length;
  const cols = matrix[0].length;
  if (rows !== cols) return false;

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > tol) return false;
    }
  }
  return true;
}

// Endpoint de verificación de estado
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'node-api' });
});

// Endpoint que calcula las estadísticas
app.post('/api/stats', (req, res) => {
  const { matrices } = req.body;
  if (!matrices || !Array.isArray(matrices)) {
    return res.status(400).json({ error: 'Se requiere el campo "matrices"' });
  }

  let values = [];
  let anyDiagonal = false;

  for (const mat of matrices) {
    if (isDiagonal(mat)) anyDiagonal = true;
    for (const row of mat) {
      for (const val of row) {
        values.push(val);
      }
    }
  }

  if (values.length === 0) {
    return res.status(400).json({ error: 'Matrices vacías' });
  }

  const sum = values.reduce((a, b) => a + b, 0);
  res.json({
    max: Math.max(...values),
    min: Math.min(...values),
    average: sum / values.length,
    sum: sum,
    is_diagonal: anyDiagonal
  });
});

app.listen(PORT, () => {
  console.log(`Node API lista en puerto ${PORT}`);
});