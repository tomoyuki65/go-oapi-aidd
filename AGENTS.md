# 全体のルール定義

## 概要

本プロジェクトはマルチエージェント構成（Codex Agents）を前提とする。

以下のエージェントが `.codex/agents` に定義されており、
すべての開発は責務分離されたTDDフローに従って実行される。

本プロジェクトは以下を前提とする：

- OpenAPI駆動開発
- TDD（Red → Green → Refactor）
- 制約ベースアーキテクチャ（DDD簡略モデル）

## エージェント構成

### pm

- 要件整理・仕様定義
- ユーザーストーリー作成
- タスク分解・優先順位決定
- 実装詳細には踏み込まない

### tester

- テスト設計（REDフェーズ）
- 失敗するテストコード作成
- 仕様の穴の検出
- API期待仕様の明文化

### implementer

- 実装（GREENフェーズ）
- テストを通す最小実装
- 設計判断は禁止（ルール準拠のみ）
- OpenAPI生成コードを前提に実装

### reviewer

- コードレビュー
- 依存ルール違反の検出
- アーキテクチャ逸脱の検出
- テストの妥当性検証
- 実装修正は行わない（再実装指示のみ）

## ワークフロー

開発は必ず以下のTDDフローに従う：

- `.codex/workflows/tdd_flow.md`

## ルール参照構造

エージェントは実装判断前に以下を参照する：

### アーキテクチャ

- `docs/rules/architecture.md`

### モジュール配置

- `docs/rules/module-classification.md`

### 依存関係

- `docs/rules/dependency-rules.md`

### OpenAPI生成ルール

- `docs/rules/oapi-codegen.md`

### ORM（Bun）

- `docs/rules/orm-bun.md`

### テスト

- `docs/rules/testing.md`

## 実装ルール

実装詳細はすべて以下に委譲する：

- `src/AGENTS.md`

※ implementer は必ずこれに従う

## ルール優先順位

矛盾が発生した場合は以下の優先順位で解決する：

1. dependency-rules.md
2. architecture.md
3. module-classification.md
4. src/AGENTS.md
5. oapi-codegen.md
6. orm-bun.md
7. testing.md

## エージェント責務原則

### 共通原則

- 各エージェントは責務外の判断を行わない
- 設計と実装は必ず分離する
- TDDフローをスキップしない
- 推測による実装禁止

## 責務分離

### pm

- 仕様決定・タスク分解
- 実装禁止

### tester

- テストのみ（RED）
- 実装禁止

### implementer

- 実装のみ（GREEN）
- 設計判断禁止
- ルール準拠のみ

### reviewer

- 検証のみ
- 実装修正禁止
- NG時は再実装指示のみ

## TDD制約（重要）

- REDなしで実装開始禁止
- REDレビューなしでGREEN禁止
- GREEN未達でレビュー禁止
- GREENレビューなしで完了禁止
- reviewer NG時は必ず再実装ループ

## 禁止事項

- 本ファイルに実装ルールを書くこと
- エージェントの責務を曖昧にすること
- TDDフローを省略すること
- reviewerが実装修正を行うこと
- 存在しないルールファイル（旧設計含む）を参照対象にしないこと
