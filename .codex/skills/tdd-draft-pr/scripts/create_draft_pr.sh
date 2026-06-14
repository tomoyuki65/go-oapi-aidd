#!/usr/bin/env bash

set -euo pipefail

# ==================================================
# Usage:
# ./create_draft_pr.sh "タイトル" "本文" "ラベル名"
# ==================================================

TITLE="${1:-}"
BODY="${2:-}"
LABELS="${3:-}"

# タイトルと本文の入力チェック
if [[ -z "$TITLE" || -z "$BODY" ]]; then
  echo "Usage: $0 \"タイトル\" \"本文\""
  exit 1
fi

# gh コマンド存在チェック
if ! command -v gh &> /dev/null; then
  echo "gh コマンドが見つかりません"
  exit 1
fi

# GitHub認証チェック
if ! gh auth status &> /dev/null; then
  echo "GitHubにログインしていません"
  echo "gh auth login を実行してください"
  exit 1
fi

echo "ドラフトPRを作成中..."

# ドラフトPR作成
if [[ -n "$LABELS" ]]; then
  gh pr create \
    --draft \
    --title "$TITLE" \
    --body-file <(printf "%s" "$BODY") \
    --label "$LABELS"
else
  gh pr create \
    --draft \
    --title "$TITLE" \
    --body-file <(printf "%s" "$BODY")
fi

echo "ドラフトPR作成完了"
