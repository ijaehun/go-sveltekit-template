#!/usr/bin/env bash
# 수동 배포 스크립트 (tar + scp). 평소엔 GitHub Actions 자동 배포를 쓰고,
# CI 없이 급하게 올려야 할 때만 사용.
set -euo pipefail

EC2_HOST="${EC2_HOST:-your.server.ip}"
EC2_USER="${EC2_USER:-ubuntu}"
KEY_FILE="${KEY_FILE:-~/.ssh/your-key.pem}"
APP_NAME="${APP_NAME:-my-app}"

echo "▶ Build artifacts..."
tar -czf "$APP_NAME.tar.gz" \
  --exclude='node_modules' \
  --exclude='.svelte-kit' \
  --exclude='build' \
  --exclude='tmp' \
  --exclude='.git' \
  --exclude='*.tar.gz' \
  .

echo "▶ Upload to $EC2_HOST..."
scp -i "$KEY_FILE" "$APP_NAME.tar.gz" "$EC2_USER@$EC2_HOST:~/"

echo "▶ Deploy on remote..."
ssh -i "$KEY_FILE" "$EC2_USER@$EC2_HOST" bash <<EOF
  set -e
  cd ~/
  rm -rf "$APP_NAME" || true
  mkdir "$APP_NAME"
  tar -xzf "$APP_NAME.tar.gz" -C "$APP_NAME/"
  cd "$APP_NAME"
  docker compose down || true
  docker compose -f docker-compose.yml up -d --build
EOF

echo "▶ Cleanup local tar..."
rm "$APP_NAME.tar.gz"

echo "✅ Deployed to http://$EC2_HOST/"
