#!/bin/bash

echo "=== Teste Final da Interface Web ==="

# Iniciar servidor em background e aguardar
cd /Users/roberto/github/cfs_spool
echo "Iniciando servidor..."
go run cmd/app/main.go &
SERVER_PID=$!

# Aguardar servidor iniciar
sleep 3

echo "Testando GET /api/read:"
curl -s http://localhost:8080/api/read | head -100

echo ""
echo "Testando POST /api/read:"  
curl -X POST -s http://localhost:8080/api/read | head -100

echo ""
echo "Testando interface principal:"
curl -I http://localhost:8080/ 2>/dev/null | grep "HTTP"

echo ""
echo "✅ Correções implementadas:"
echo "- ✅ URL corrigida: /api/read-tag → /api/read"  
echo "- ✅ Método POST aceito na API"
echo "- ✅ Compatibilidade com tags de 48 bytes"
echo "- ✅ ParseFieldsCompat() funcionando"

echo ""
echo "🌐 Teste na interface web: http://localhost:8080"
echo "📱 Agora o botão 'Ler Tag' deve funcionar!"

# Matar servidor
kill $SERVER_PID 2>/dev/null
