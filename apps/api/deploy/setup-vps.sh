#!/bin/bash
# =============================================================================
# Go Framework Backend — VPS Setup (Docker-based)
# =============================================================================
# Jalankan di VPS sebagai root:
#   sudo bash setup-vps.sh
# =============================================================================

set -e

echo "========================================="
echo "  Go Framework Backend — VPS Setup"
echo "========================================="

# 1. Update system
echo "[1/4] Updating system..."
apt update && apt upgrade -y

# 2. Install Docker
echo "[2/4] Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
fi
docker --version
docker compose version

# 3. Install Nginx + Certbot
echo "[3/4] Installing Nginx & Certbot..."
apt install -y nginx certbot python3-certbot-nginx git
systemctl enable nginx

# 4. Clone repo
echo "[4/4] Setting up project..."
read -p "Nama project (default: myapp): " PROJECT_NAME
PROJECT_NAME="${PROJECT_NAME:-myapp}"
PROJECT_DIR="/var/www/$PROJECT_NAME"

if [ ! -d "$PROJECT_DIR" ]; then
    read -p "GitHub repo URL (e.g. https://github.com/username/repo.git): " REPO_URL
    git clone -b staging "$REPO_URL" "$PROJECT_DIR"
    cd "$PROJECT_DIR"
    cp .env.example .env
    echo ""
    echo ">> Edit .env sebelum lanjut:"
    echo "   nano $PROJECT_DIR/.env"
    echo ""
    echo ">> Setelah edit .env, jalankan:"
    echo "   cd $PROJECT_DIR && docker compose up -d"
    echo ""
    echo ">> Setup Nginx + SSL:"
    echo "   cp deploy/nginx.conf /etc/nginx/sites-available/$PROJECT_NAME"
    echo "   ln -s /etc/nginx/sites-available/$PROJECT_NAME /etc/nginx/sites-enabled/"
    echo "   nginx -t && systemctl reload nginx"
    echo "   certbot --nginx -d api.yourdomain.com"
else
    echo "Project sudah ada di $PROJECT_DIR"
    cd "$PROJECT_DIR"
    git pull origin staging
    docker compose up -d --build
fi

echo ""
echo "========================================="
echo "  Setup selesai!"
echo "========================================="
