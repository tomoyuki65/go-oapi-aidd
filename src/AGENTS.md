# Goの実装ルール

このファイルをGoの実装における「実行ルールの入口」として扱う。  
詳細な設計ルールはすべて `docs/rules` 配下を参照すること。

Go実装時は、まずこのファイルを読み、実装プロセスを確認すること。
そのうえで、判断に必要な詳細ルールを以下から参照する。

1. `docs/rules/dependency-rules.md`
2. `docs/rules/architecture.md`
3. `docs/rules/module-classification.md`
4. `docs/rules/oapi-codegen.md`
5. `docs/rules/orm-bun.md`
6. `docs/rules/testing.md`

矛盾が発生した場合は、リポジトリ直下の `AGENTS.md` の「ルール優先順位」に従うこと。

## 1. 基本原則

- 必ずTDD（テスト駆動開発）で実装する
- Red → Green → Refactor を厳守する
- 既存構造を勝手に変更しない
- ルールに存在しない構造は作らない
- すべての判断は `docs/rules` を最優先とする
- 推測で設計判断を行わない

## 2. 実装プロセス（必須手順）

新規APIの作成・既存APIの変更は必ず以下の順序で行う。
API仕様の確定とコード生成はREDテスト作成より前に完了させる。

### 1. モジュール分類

まず `docs/rules/module-classification.md` を参照し、対象機能を分類する。

分類先：

- core
- supporting
- generic
- shared

配置先は分類結果に従うこと。

### 2. OpenAPI定義を修正する

APIの追加・変更を行う場合は、実装より先にOpenAPI定義を更新する。

API仕様を先に確定させること。

対象：

- `src/openapi/openapi.yaml`
- `src/openapi/paths/*`
- `src/openapi/components/schemas/*`

### 3. コード生成を実行する

OpenAPI定義を更新した場合は、必ず以下のコマンドを実行し、コード生成を行う。

```bash
make generate
```

生成コードは以下に出力される。

- `src/internal/presentation/gen`

生成コードは手動編集してはならない。

禁止事項：

- `make generate` 実行前に生成コードを編集してはならない
- `src/internal/presentation/gen` 配下を手動編集してはならない

### 4. テストコードを先に作成する（Red）

実装前に必ず失敗するテストを書く。

テストルールは `docs/rules/testing.md` に従うこと。

### 5. 実装する（Green）

テストを通すための最小実装を行う。

以下のルールを厳守すること。

- `docs/rules/architecture.md`
- `docs/rules/dependency-rules.md`
- `docs/rules/module-classification.md`
- `docs/rules/oapi-codegen.md`
- `docs/rules/orm-bun.md`

### 6. リファクタリングする（Refactor）

テスト成功後にのみリファクタリングを行う。

- 振る舞いを変更しない
- テストを壊さない
- 構造のみ改善する

## 3. 判断ルール

判断に迷った場合は以下を参照する。

- `docs/rules/dependency-rules.md`
- `docs/rules/architecture.md`
- `docs/rules/module-classification.md`
- 既存コード

それでも判断できない場合は、`docs/rules/module-classification.md` の分類基準を再確認すること。

推測による配置判断・設計判断は禁止する。

## 4. 禁止事項

以下を禁止する。

- OpenAPI定義を変更せずにAPIを追加すること
- OpenAPI定義より先にAPI実装を開始すること
- 生成コードを手動編集すること
- `docs/rules` に存在しない構造を追加すること
- 推測でモジュール分類を行うこと
- 推測で依存関係を追加すること
- テストを書かずに実装を開始すること
- Red を経ずに Green を行うこと
