# ORM「Bun」に関するルール定義

このファイルは ORM「Bun」の利用ルールを定義する。

Bun の利用可能場所・利用禁止場所・マイグレーション運用を定義し、DB アクセスの責務を明確化する。

## 1. 基本方針

本プロジェクトの DB アクセスは Bun を使用する。

SQL アクセスは Bun を経由して実装すること。

Bun の利用場所は明示的に制限される。

モジュール分類ごとに異なる利用ルールを適用する。

- core：Repository パターン
- supporting：Transaction Script
- generic：Active Record

詳細なモジュール分類は `module-classification.md` を参照すること。

## 2. Bun の利用可能場所

Bun は以下でのみ使用できる。

- `internal/core/*/infrastructure/repository`
- `internal/core/*/domain/*_repository.go` の repository interface 引数型（`bun.IDB` のみ）
- `internal/supporting`
- `internal/generic`
- `internal/infrastructure/database`
- `internal/di`

## 3. Bun の利用禁止場所

以下では Bun を使用してはならない。

- `internal/core/*/domain`
  - 例外：repository interface の引数型として `bun.IDB` を使う場合のみ許可
- `internal/core/*/usecase`
- `internal/core/*/infrastructure/external`
- `internal/shared`
- `internal/presentation`

## 4. core での利用ルール

core は DDD ベースの構造を採用する。

### domain

domain は業務ルールのみを扱う。

禁止事項：

- Bun の利用
- SQL の記述
- Bun のタグ利用
- Bun の型への依存
  - 例外：repository interface の引数型として `bun.IDB` を使う場合のみ許可
- DB アクセス

許可事項：

- repository interface では、通常の `*bun.DB` と `bun.Tx` の両方を受け取れるようにする目的で、引数型に `bun.IDB` を使ってよい。
- `bun.IDB` は DB 実行主体を表す型としてのみ扱い、domain 内で query builder の生成、SQL記述、DBアクセスを行ってはならない。

例：

```go
type MemberRepository interface {
    FindByID(ctx context.Context, db bun.IDB, id string) (*Member, error)
}
```

### usecase

usecase は業務フローのみを扱う。

禁止事項：

- Bun の利用
- SQL の記述
- DB への直接アクセス

### infrastructure/repository

DB アクセスは repository が担当する。

Bun は repository 実装内で利用する。

対象：

- command（書き込み）
- query（読み取り）

Repository は domain に定義されたインターフェースを実装する。

### infrastructure/external

外部サービス連携を担当する。

禁止事項：

- Bun の利用
- DB アクセス

---

## 5. supporting での利用ルール

supporting は Transaction Script を採用する。

DB アクセスが必要な場合は `service.go` で Bun を利用してよい。

ただし以下を守ること。

- 単一業務フローとして完結させる
- 業務処理と DB 処理を過度に混在させない
- 不要に複雑な SQL を記述しない
- ドメインモデルを導入しない

## 6. generic での利用ルール

generic は Active Record を採用する。

単純 CRUD については Bun を直接利用してよい。

対象：

- Create
- Read
- Update
- Delete

禁止事項：

- 複雑な業務ルールの実装
- ドメインロジックの実装
- 複数ステップの整合性制御
- 長大なトランザクション処理

generic は技術的処理のみを扱う。

## 7. スキーマ定義

Bun用のスキーマ定義は `internal/infrastructure/database/schema` 配下に配置する。

責務：

- テーブル定義とのマッピング
- Bun タグの定義

禁止事項：

- 業務ロジックの実装
- バリデーションの実装
- ドメインモデルとの兼用

スキーマ定義は永続化専用モデルとして扱う。

## 8. 接続定義

Bunの接続設定は `internal/infrastructure/database/bun.go` を利用する。

責務：

- DB 接続生成
- 接続設定
- Bun 初期化

禁止事項：

- 業務ロジックの実装
- Repository の生成

## 9. DI コンテナ

DI コンテナは Bun を利用してよい。

配置：

- `internal/di/container.go`

責務：

- Bun 接続の注入
- Repository の生成
- Service の生成
- Handler の生成
- 依存関係の組み立て

禁止事項：

- 業務ロジックの実装
- SQL の記述

## 10. マイグレーション

マイグレーション関連ファイルは `internal/infrastructure/migration` 配下に配置する。

マイグレーションは必ず Bun のマイグレーション機能を利用する。

### ファイル作成ルール

マイグレーションファイルは手動作成してはならない。

必ず以下のコマンドを実行する。

```bash
docker compose exec api go run cmd/migrate/main.go create_sql [ファイル名]
```

例：

```bash
docker compose exec api go run cmd/migrate/main.go create_sql create_users_table
```

実行すると以下のような2つのファイルが生成される。

```text
internal/infrastructure/migrations/20260119023405_create_users_table.up.sql
internal/infrastructure/migrations/20260119023405_create_users_table.down.sql
```

- `*.up.sql`：マイグレーション実行用
- `*.down.sql`：ロールバック用

### 命名規則

ファイル名は以下の形式とする。

```text
[操作]_[対象]
```

例：

- `create_users_table`
- `add_index_to_orders`
- `drop_legacy_columns`

タイムスタンプは自動付与されるため手動指定しない。

### 初期化

マイグレーション管理テーブルを作成する場合は以下を実行する。

```bash
docker compose exec api go run cmd/migrate/main.go init
```

テスト用 DB に対して実行する場合：

```bash
docker compose exec api env ENV=testing go run cmd/migrate/main.go init
```

### 状態確認

現在のマイグレーション状態を確認する場合は以下を実行する。

```bash
docker compose exec api go run cmd/migrate/main.go status
```

### マイグレーション実行

マイグレーションを適用する場合は以下を実行する。

```bash
docker compose exec api go run cmd/migrate/main.go migrate
```

一度に実行された SQL は同一グループとして管理される。

### ロールバック

直前のマイグレーショングループをロールバックする場合は以下を実行する。

```bash
docker compose exec api go run cmd/migrate/main.go rollback
```

ロールバックはファイル単位ではなくグループ単位で実行される。

### マイグレーション作成ルール

- `up.sql` と `down.sql` は必ず対で作成する
- `down.sql` を省略してはならない
- 既存の適用済みマイグレーションを書き換えてはならない
- 過去のマイグレーションを削除してはならない
- スキーマ変更は新規マイグレーションとして追加する

## 11. 禁止事項

以下を禁止する。

- domain で Bun を使用すること
- domain で SQL を記述すること
- usecase で Bun を使用すること
- usecase で SQL を記述すること
- presentation で Bun を使用すること
- shared で Bun を使用すること
- DI コンテナに業務ロジックを書くこと
- Bun のスキーマを業務モデルとして扱うこと
- マイグレーションファイルを手動作成すること
- 適用済みマイグレーションを書き換えること

## 12. 判断基準

DB アクセスが必要な場合は以下で判断する。

### core

Repository に実装する。

### supporting

`service.go` に実装する。

### generic

Active Record として実装する。

判断に迷った場合は `module-classification.md` を参照し、対象機能の分類を先に決定すること。
