#!/bin/bash
# Script para gerar ícones .icns para macOS a partir do SVG

# Verifica se o ImageMagick está instalado
if ! command -v convert &> /dev/null; then
    echo "ImageMagick não encontrado. Instalando via Homebrew..."
    brew install imagemagick
fi

# Cria diretório temporário para os PNGs
mkdir -p assets/icons/iconset.iconset

# Tamanhos necessários para .icns
sizes=(16 32 64 128 256 512 1024)

echo "Gerando ícones PNG a partir do SVG..."

for size in "${sizes[@]}"; do
    echo "Gerando ícone ${size}x${size}..."
    convert -background none assets/icons/cfs-spool.svg -resize ${size}x${size} assets/icons/iconset.iconset/icon_${size}x${size}.png
    
    # Para retina displays (2x)
    if [ $size -le 512 ]; then
        double_size=$((size * 2))
        convert -background none assets/icons/cfs-spool.svg -resize ${double_size}x${double_size} assets/icons/iconset.iconset/icon_${size}x${size}@2x.png
    fi
done

# Gera o arquivo .icns usando iconutil
echo "Gerando arquivo .icns..."
iconutil -c icns assets/icons/iconset.iconset -o assets/icons/cfs-spool.icns

# Limpa arquivos temporários
rm -rf assets/icons/iconset.iconset

echo "Ícone .icns criado em: assets/icons/cfs-spool.icns"
