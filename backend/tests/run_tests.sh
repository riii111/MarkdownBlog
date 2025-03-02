#!/bin/bash

# テスト実行スクリプト

# 色の設定
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# テストフォルダのベースパス
TEST_BASE_PATH="./tests/e2e"

# テスト結果を格納する配列
declare -a TEST_RESULTS
declare -a TEST_FOLDERS

echo -e "${YELLOW}Markdown Blog API E2Eテスト実行${NC}"
echo "========================================"

# テストフォルダを取得
TEST_DIRS=$(find "$TEST_BASE_PATH" -type d -name "*_test.go" -prune -o -type d -not -path "*/\.*" -print | sort)

if [ -z "$TEST_DIRS" ]; then
  echo -e "${RED}テストディレクトリの取得に失敗しました${NC}"
  exit 1
fi

# テスト実行
for dir in $TEST_DIRS; do
  # ベースパスからの相対パスを取得
  if [ "$dir" != "$TEST_BASE_PATH" ] && [ -d "$dir" ] && [ "$(ls -A "$dir" 2>/dev/null)" ]; then
    # テストファイルが存在するか確認
    if ls "$dir"/*_test.go 1> /dev/null 2>&1; then
      FOLDER_NAME=$(basename "$dir")
      echo -e "${BLUE}テストフォルダ: ${YELLOW}$FOLDER_NAME${NC} 実行中..."
      
      # テスト実行
      go test -v "$dir/..."
      RESULT=$?
      
      # 結果を配列に格納
      TEST_RESULTS+=("$RESULT")
      TEST_FOLDERS+=("$FOLDER_NAME")
      
      if [ $RESULT -eq 0 ]; then
        echo -e "${GREEN}$FOLDER_NAME テスト成功${NC}"
      else
        echo -e "${RED}$FOLDER_NAME テスト失敗${NC}"
      fi
      
      echo "----------------------------------------"
    fi
  fi
done

echo "========================================"

# 全体の結果を表示
ALL_PASSED=true

for i in "${!TEST_RESULTS[@]}"; do
  if [ ${TEST_RESULTS[$i]} -ne 0 ]; then
    ALL_PASSED=false
    echo -e "${RED}${TEST_FOLDERS[$i]} テスト失敗${NC}"
  fi
done

if $ALL_PASSED; then
  echo -e "${GREEN}全てのテストが成功しました${NC}"
  exit 0
else
  echo -e "${RED}テストに失敗があります${NC}"
  exit 1
fi
