# Goの実装ルール定義

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
API仕様の確定とコード生成はRedテスト作成より前に完了させる。

### 1. モジュール分類

まず `docs/rules/module-classification.md` を参照し、対象機能を分類する。

分類先：

- core
- supporting
- generic
- shared

配置先は分類結果に従うこと。

---

### 2. OpenAPI定義を修正する

APIの追加・変更を行う場合は、実装より先にOpenAPI定義を更新する。

API仕様を先に確定させること。

対象：

- `src/openapi/openapi.yaml`
- `src/openapi/paths/*`
- `src/openapi/components/schemas/*`

---

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

---

### 4. テストコードを先に作成する（Red）

実装前に必ず失敗するテストを書く。

テストルールは docs/rules/testing.md に従うこと。

テスト追加後は、対象テストを実行し、期待どおり失敗することを確認する。

実行コマンド：

・ユニットを追加した場合

```bash
make test-unit
```

・インテグレーションテストを追加した場合

```bash
make test-integration
```

・e2eテストを追加した場合

```bash
make test-e2e
```

失敗を確認せずに実装を開始してはならない。

---

### 5. 実装する（Green）

テストを通すための最小実装を行う。

以下のルールを厳守すること。

- `docs/rules/architecture.md`
- `docs/rules/dependency-rules.md`
- `docs/rules/module-classification.md`
- `docs/rules/oapi-codegen.md`
- `docs/rules/orm-bun.md`

実装後は、追加・変更したすべてのテストが成功することを確認する。

実行コマンド：

```bash
make test
```

個別実行する場合：

・ユニットテスト実行

```bash
make test-unit
```

・インテグレーションテスト実行

```bash
make test-integration
```

・e2eテスト実行（存在する場合のみ）

```bash
make test-e2e
```

Green未達の状態でリファクタリングを行ってはならない。

---

### 6. リファクタリングする（Refactor）

テスト成功後にのみリファクタリングを行う。

- 振る舞いを変更しない
- テストを壊さない
- 構造のみ改善する

---

### 7. フォーマット・静的解析を実行する

コードの追加・修正後は、必ずフォーマットおよび静的コード解析を実行する。

実装完了条件は、以下のコマンドがすべて成功すること。

・フォーマット修正

```bash
docker compose run --rm api golangci-lint fmt -v ./...
```

・静的コード解析

```bash
docker compose run --rm api golangci-lint run -v ./...
```

違反が検出された場合は、すべて解消してからレビュー依頼を行うこと。

生成コードを含め、リポジトリ全体がチェック対象となる。

---

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

- OpenAPI定義を変更せずにAPIを追加・変更すること
- OpenAPI定義より先にAPI実装を開始すること
- make generate を実行せずに生成コードを利用すること
- 生成コードを手動編集すること
- docs/rules に存在しない構造を追加すること
- 推測でモジュール分類を行うこと
- 推測で依存関係を追加すること
- テストを書かずに実装を開始すること
- Redを経ずにGreenを行うこと
- Redを確認せずに実装を開始すること
- Greenを確認せずにリファクタリングを行うこと
- フォーマットを実行せずにレビュー依頼すること
- 静的コード解析を実行せずにレビュー依頼すること
- テストが失敗した状態でレビュー依頼すること
- golangci-lint run でエラーが残った状態でコミットすること
- テストが失敗した状態でコミットすること
