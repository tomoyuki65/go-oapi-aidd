#!/usr/bin/env bash

set -euo pipefail

# ==================================================
# Usage:
# ./create_issue.sh "タイトル" "本文" "ラベル名"
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

echo "Issueを作成中..."

# Issue作成
if [[ -n "$LABELS" ]]; then
  gh issue create \
    --title "$TITLE" \
    --body-file <(printf "%s" "$BODY") \
    --label "$LABELS"
else
  gh issue create \
    --title "$TITLE" \
    --body-file <(printf "%s" "$BODY")
fi

echo "Issue作成完了"
