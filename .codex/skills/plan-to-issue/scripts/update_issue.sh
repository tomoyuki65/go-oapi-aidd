#!/usr/bin/env bash

set -euo pipefail

# ==================================================
# Usage:
# ./update_issue.sh "Issue番号" "本文" "ラベル名"
# ==================================================

ISSUE_NUMBER="${1:-}"
BODY="${2:-}"
LABELS="${3:-}"

# Issue番号と本文の入力チェック
if [[ -z "$ISSUE_NUMBER" || -z "$BODY" ]]; then
  echo "Usage: $0 \"Issue番号\" \"本文\""
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

echo "Issueを更新中..."

# Issue更新
if [[ -n "$LABELS" ]]; then
  gh issue edit "$ISSUE_NUMBER" \
    --body-file <(printf "%s" "$BODY") \
    --add-label "$LABELS"
else
  gh issue edit "$ISSUE_NUMBER" \
    --body-file <(printf "%s" "$BODY")
fi

echo "Issue更新完了"
