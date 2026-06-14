#!/usr/bin/env bash

set -euo pipefail

# ==================================================
# Usage:
# ./add_pr_comment.sh "PR番号" "追加コメント内容"
# ==================================================

PR_NUMBER="${1:-}"
BODY="${2:-}"

# PR番号と追加コメント内容の入力チェック
if [[ -z "$PR_NUMBER" || -z "$BODY" ]]; then
  echo "Usage: $0 \"PR番号\" \"追加コメント内容\""
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

echo "PRにコメントを投稿中..."

# PRにコメントを投稿
gh pr comment "$PR_NUMBER" --body-file <(printf "%s" "$BODY")

echo "PRにコメントを投稿しました。"
