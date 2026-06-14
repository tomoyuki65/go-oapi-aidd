#!/usr/bin/env bash

set -euo pipefail

# ==================================================
# Usage:
# ./mark_pr_ready.sh "PR番号" "追加コメント内容"
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

echo "ドラフトPRにコメントを投稿し、ready for review に変更中..."

# ドラフトPRにコメントを投稿
gh pr comment "$PR_NUMBER" --body-file <(printf "%s" "$BODY")

# ドラフトPRを ready for review 状態に変更
gh pr ready "${PR_NUMBER}"

echo "PRを更新しました。"
