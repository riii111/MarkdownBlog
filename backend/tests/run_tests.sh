#!/bin/bash

# テスト実行スクリプト

# 色の設定
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Markdown Blog API E2Eテスト実行${NC}"
echo "========================================"

# 認証系テスト
echo -e "${YELLOW}認証系APIテスト実行中...${NC}"
go test -v ./tests/e2e/auth/...
AUTH_RESULT=$?

if [ $AUTH_RESULT -eq 0 ]; then
  echo -e "${GREEN}認証系APIテスト成功${NC}"
else
  echo -e "${RED}認証系APIテスト失敗${NC}"
fi

echo "========================================"

# 全体の結果
if [ $AUTH_RESULT -eq 0 ]; then
  echo -e "${GREEN}全てのテストが成功しました${NC}"
  exit 0
else
  echo -e "${RED}テストに失敗があります${NC}"
  exit 1
fi
