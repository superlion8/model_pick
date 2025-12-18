const fs = require('fs');
const path = require('path');

// 读取生成的 models_data.json
const modelsDataPath = path.join(__dirname, '..', 'models_data.json');
let modelsData = [];

try {
  const data = fs.readFileSync(modelsDataPath, 'utf8');
  modelsData = JSON.parse(data);
} catch (error) {
  console.error('Error reading models_data.json:', error);
}

module.exports = (req, res) => {
  // 设置CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');

  if (req.method === 'OPTIONS') {
    res.status(200).end();
    return;
  }

  // 返回模特数据
  res.status(200).json(modelsData);
};
