#!/bin/bash

# Script para verificar e testar o sistema de tagueamento automático
set -e

# Cores para saída
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Teste do Sistema de Tagueamento Automático${NC}"
echo "=================================================="

# Verificar branch atual
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo -e "${BLUE}📋 Branch atual: ${YELLOW}${CURRENT_BRANCH}${NC}"

# Buscar todas as tags
git fetch --tags --force
echo -e "${BLUE}📋 Tags atuais:${NC}"
git tag -l "v*" | sort -V

# Obter última tag
LATEST_TAG=$(git tag -l "v[0-9]*.[0-9]*.[0-9]*" | sort -V | tail -n 1)
if [ -z "$LATEST_TAG" ]; then
    LATEST_TAG="v0.0.0"
    echo -e "${YELLOW}⚠️  Nenhuma tag encontrada, iniciando com ${LATEST_TAG}${NC}"
else
    echo -e "${BLUE}📋 Última tag: ${YELLOW}${LATEST_TAG}${NC}"
fi

# Extrair componentes
VERSION=${LATEST_TAG#v}
MAJOR=$(echo $VERSION | cut -d. -f1)
MINOR=$(echo $VERSION | cut -d. -f2)
PATCH=$(echo $VERSION | cut -d. -f3)

echo -e "${BLUE}📋 Componentes da versão:${NC}"
echo -e "   MAJOR: ${YELLOW}${MAJOR}${NC}"
echo -e "   MINOR: ${YELLOW}${MINOR}${NC}"
echo -e "   PATCH: ${YELLOW}${PATCH}${NC}"

# Simulação de incrementos
echo -e "\n${BLUE}📋 Simulação de incrementos:${NC}"

# Patch
NEW_PATCH=$((PATCH + 1))
echo -e "   Patch: ${LATEST_TAG} → ${GREEN}v${MAJOR}.${MINOR}.${NEW_PATCH}${NC}"

# Minor
NEW_MINOR=$((MINOR + 1))
echo -e "   Minor: ${LATEST_TAG} → ${GREEN}v${MAJOR}.${NEW_MINOR}.0${NC}"

# Major
NEW_MAJOR=$((MAJOR + 1))
echo -e "   Major: ${LATEST_TAG} → ${GREEN}v${NEW_MAJOR}.0.0${NC}"

echo -e "\n${BLUE}🧪 Passos para testar o workflow:${NC}"
echo -e "1. Certifique-se que você está na branch ${YELLOW}main${NC}"
echo -e "   ${GREEN}git checkout main${NC}"
echo -e "2. Faça uma alteração simples"
echo -e "3. Commit com o tipo de incremento desejado:"
echo -e "   ${GREEN}git commit -m \"Mensagem #patch\"${NC} (para incremento de patch)"
echo -e "   ${GREEN}git commit -m \"Mensagem #minor\"${NC} (para incremento de minor)"
echo -e "   ${GREEN}git commit -m \"Mensagem #major\"${NC} (para incremento de major)"
echo -e "4. Envie para o repositório remoto:"
echo -e "   ${GREEN}git push origin main${NC}"
echo -e "5. Verifique se o workflow foi iniciado em:"
echo -e "   ${BLUE}https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\).git/\1/')/actions${NC}"

echo -e "\n${BLUE}🧪 Teste Manual:${NC}"
echo -e "Você também pode acionar o workflow manualmente via interface do GitHub:"
echo -e "1. Acesse Actions → Auto Tag → Run workflow"
echo -e "2. Selecione o tipo de incremento (patch, minor, major)"
echo -e "3. Clique em Run workflow"

echo -e "\n${YELLOW}📝 Nota: Se o workflow falhar, verifique os logs na interface do GitHub Actions${NC}"
