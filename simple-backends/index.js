const express = require('express');
const app = express();
const port = 8009;

// Định nghĩa route cho /health
app.get('/health', (req, res) => {
  res.send('Server is healthy');
});

// Lắng nghe các yêu cầu trên cổng đã chỉ định
app.listen(port, () => {
  console.log(`Server is running at http://localhost:${port}`);
});
