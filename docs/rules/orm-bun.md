# ORM「Bun」に関するルール定義

このファイルは ORM「Bun」の利用ルールを定義する。

Bunの利用可能場所と利用禁止場所を定義し、DBアクセスの責務を明確化する。

## 1. 基本方針

本プロジェクトのDBアクセスは Bun を使用する。

SQLアクセスは Bun を経由して実装すること。

Bunの利用場所は明示的に制限される。

## 2. Bunの利用可能場所

Bunは以下でのみ使用できる。

- `internal/core/*/infrastructure/repository`
- `internal/supporting`
- `internal/generic`
- `internal/infrastructure/database`
- `internal/di`

## 3. Bunの利用禁止場所

以下では Bun を使用してはならない。

- `internal/core/*/domain`
- `internal/core/*/usecase`
- `internal/core/*/infrastructure/external`
- `internal/shared`
- `internal/presentation`

## 4. coreでの利用ルール

coreはDDD構造を採用する。

### domain

domainは業務ルールのみを扱う。

禁止事項：

- Bunの利用
- SQLの記述
- Bunのタグ利用
- Bunの型への依存

### usecase

usecaseは業務フローのみを扱う。

禁止事項：

- Bunの利用
- SQLの記述
- DBへの直接アクセス

### infrastructure/repository

DBアクセスは repository が担当する。

Bunは repository 実装内で利用する。

対象：

- command（書き込み）
- query（読み取り）

### infrastructure/external

外部サービス連携を担当する。

禁止事項：

- Bunの利用
- DBアクセス

## 5. supportingでの利用ルール

supportingは transaction script を採用する。

DBアクセスが必要な場合は service.go で Bun を利用してよい。

ただし以下を守ること。

- 業務処理とDB処理を過度に混在させない
- 不要な複雑SQLを記述しない
- 単一業務フローとして完結させる

## 6. genericでの利用ルール

genericは Active Record を採用する。

単純CRUDについては Bun を直接利用してよい。

対象：

- 作成（Create）
- 取得（Read）
- 更新（Update）
- 削除（Delete）

禁止事項：

- 複雑な業務ルールの実装
- ドメインロジックの実装
- トランザクション整合性を伴う業務処理

## 7. スキーマ定義

Bun用のスキーマ定義は以下に配置する。

- `internal/infrastructure/database/schema`

スキーマ定義はDBマッピングの責務のみを持つ。

業務ロジックを含めてはならない。

## 8. 接続定義

Bunの接続設定は以下に配置する。

- `internal/infrastructure/database/bun.go`

接続生成・設定のみを担当する。

## 9. DIコンテナ

DIコンテナは Bun を利用してよい。

配置：

- `internal/di/container.go`

責務：

- Bun接続の注入
- Repositoryの生成
- Serviceの生成
- Handlerの生成
- 依存関係の組み立て

DIコンテナに業務ロジックを書いてはならない。

## 10. 禁止事項

以下を禁止する。

- domainで Bun を使用すること
- domainで SQL を記述すること
- usecaseで Bun を使用すること
- usecaseで SQL を記述すること
- presentationで Bun を使用すること
- sharedで Bun を使用すること
- DIコンテナに業務ロジックを書くこと
- Bunのスキーマを業務モデルとして扱うこと

## 11. 判断基準

DBアクセスが必要な場合は以下で判断する。

### core

- repository に実装する

### supporting

- service.go に実装する

### generic

- Active Record として実装する

判断に迷った場合は `module-classification.md` を参照し、対象機能の分類を先に決定すること。
