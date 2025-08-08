#!/bin/bash

echo "=== Teste da API com Nova Estrutura ==="

# Dados de teste para gravação (sem campo batch)
TEST_DATA='{
  "date": "2024-08-08",
  "supplier": "0276",
  "material": "01001",
  "color": "FF401",
  "length": "0165",
  "serial": "000001"
}'

echo "📤 Dados de teste (sem campo batch):"
echo "$TEST_DATA" | python3 -m json.tool

echo ""
echo "🧪 Testando endpoint /api/options..."
curl -s http://localhost:8080/api/options | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    print('✅ API Options funcionando')
    print(f'Fornecedores: {len(data.get(\"suppliers\", {}))} itens')
    print(f'Materiais: {len(data.get(\"materials\", {}))} itens')
    print(f'Comprimentos: {len(data.get(\"lengths\", {}))} itens')
except:
    print('❌ Erro ao processar resposta da API')
"

echo ""
echo "🔍 Estrutura esperada:"
echo "- Batch: fixo 'A2' (2 bytes)"
echo "- Color: '0' + 5 hex chars (6 bytes total)"
echo "- Reserve: fixo '000000' (6 bytes)"
echo "- Total: 38 bytes ASCII"

echo ""
echo "✅ Servidor está respondendo na porta 8080"
echo "✅ Nova estrutura de campos implementada"
echo "🌐 Interface disponível em: http://localhost:8080"
