# Go言語によるAI駆動開発のサンプル
  
Go言語（Golang）のoapi-codegen（chi）とBunを利用し、
AI駆動開発でDDDベースのAPIを開発する方法についてのサンプルです。
  
<br/>
  
## 動作要件

本プロジェクトの動作および開発には、以下の環境・ツールを使用します。

### 実行環境

- Go: 1.26.3
- air-verse/air: v1.65.3
- golangci-lint: v2.12.2
- go.uber.org/mock/mockgen: v0.6.0
- oapi-codegen: v2.7.0
  - chi: v5.0.3
- PostgreSQL: 18.3

### 開発ツール

- Codexアプリを活用したAI駆動開発
  - 公式ドキュメント「https://developers.openai.com/codex/app 」
  
<br/>
  
## AIツールとハーネス設計

AIツールはOpenAIの`Codex`を利用する例としており、ハーネス設計もしています。

### 1. コーデックスの設定

- `.codex/config.toml`

### 2. コーデックスのサンドボックス外禁止コマンド設定

- `.codex/rules/default.rules`

### 3. マルチエージェント設定

- `.codex/agents/pm.toml` （プロダクトマネージャー・指揮者）
- `.codex/agents/tester.toml` （テスター）
- `.codex/agents/implementer.toml` （実装者）
- `.codex/agents/reviewer.toml` （レビュワー）

### 4. ワークフロー設定

- `.codex/workflows/tdd_flow.md` （TDD開発フロー）

### 5. Agent Skills

- `.codex/skills/plan-to-issue`
  - プランモードで作成した開発計画をGitHubのIssue用のフォーマットへ変換して自動登録する。必要に応じて登録後のIssue内容の修正もする。
- `.codex/skills/auto-commit`
  - 修正したコードをステージングしたうえで差分を解析し、適切なコミットメッセージを生成してgit commitまでを自動で実行する
- `.codex/skills/tdd-draft-pr`
  - Issueと追加したREDテストをもとに、テストレビュー用のドラフトPRを作成する
- `.codex/skills/pr-sync-comment`
  - 直前のコミット内容をもとにPRコメントを生成し、ブランチをpushした上でPRにコメントを追加する
- `.codex/skills/tdd-ready-pr`
  - レビュー済みのドラフトPRに実装レビューコメントを追加し、PRをReady for Reviewへ変更する

### 6. プロジェクト用

- `AGENTS.md` （全体のルール定義）
- `src/AGENTS.md` （Goの実装ルール定義）
- docs/rules （各種ルールの詳細定義）
  - `dependency-rules.md` （依存関係のルール定義）
  - `architecture.md` （アーキテクチャ設計のルール定義）
  - `module-classification.md` （モジュール分類のルール定義）
  - `oapi-codegen.md` （oapi-codegenに関するルール定義）
  - `orm-bun.md` （ORM「Bun」に関するルール定義）
  - `testing.md` （テストのルール定義）
  
<br />
  
## ディレクトリ構成

今回はOpenAPI定義からGoコードを生成し、TDD（Test-Driven Development）を前提として開発を進めます。

アーキテクチャにはDDD（ドメイン駆動設計）の考え方を取り入れていますが、アプリケーション全体へ一律に適用することは想定していません。

DDDを全面的に適用すると、実装や運用の複雑性が増し、必ずしも費用対効果が見合わないケースも多いため、システムを以下の3つの業務領域に分類します。

- 中核の業務領域
- 補完的な業務領域
- 一般的な業務領域

このうち、ビジネス上の競争優位性に直結する「中核の業務領域」に対してのみ、ドメイン単位でモジュール化したDDD構成を適用します。

各モジュールは必要に応じて「domain」、「usecase」、「infrastructure」などのレイヤーで構成し、境界づけられたコンテキスト（Bounded Context）を意識した設計とします。

一方、「補完的な業務領域」や「一般的な業務領域」については、過度な抽象化を避け、シンプルな設計を採用します。

以上の方針を踏まえ、今回のディレクトリ構成は以下のように設計しています。  
  
```
/go-oapi-adii
 |
 └── /src
      |
      ├── /cmd/migrate/main.go（DBのマイグレーション用スクリプト）
      |
      ├── /internal
      |    |
      |    ├── /core（中核の業務領域）
      |    |    |
      |    |    └── /[domain_name]（ドメインモジュール）
      |    |         |
      |    |         ├── /domain（ドメイン）
      |    |         |   |
      |    |         |   ├── [entity_name].go（エンティティ）
      |    |         |   |
      |    |         |   ├── [valueobject_name].go（値オブジェクト）
      |    |         |   |
      |    |         |   ├── [domain_service_name]_service.go（ドメインサービス）
      |    |         |   |
      |    |         |   ├── [repository_name]_repository.go（リポジトリのインターフェース）
      |    |         |   |
      |    |         |   └── [gateway_name]_gateway.go（外部サービス用のインターフェース）
      |    |         |
      |    |         ├── /usecase（ユースケース層）
      |    |         |
      |    |         └── /infrastructure（ドメイン用のインフラストラクチャ層）
      |    |              |
      |    |              ├── /repository（リポジトリの実装）
      |    |              |    |
      |    |              |    ├── /command（書き込み）
      |    |              |    |
      |    |              |    └── /query（読み込み）
      |    |              |
      |    |              └── /external（外部サービスの実装）
      |    |
      |    ├── /supporting（補完的なの業務領域）
      |    |    |
      |    |    └── /[supporting_name]（サービス層）
      |    |         |
      |    |         └── service.go
      |    |
      |    ├── /generic（一般的の業務領域）
      |    |
      |    ├── /shared（横断関心）
      |    |
      |    ├── /di（DIコンテナ層）
      |    |    |
      |    |    └── container.go
      |    |
      |    ├── /infrastructure（共通インフラストラクチャ層）
      |    |    |
      |    |    ├── /database（データベース設定）
      |    |    |    |
      |    |    |    ├── /schema（Bun用のスキーマ定義）
      |    |    |    |
      |    |    |    └── bun.go （ORM「Bun」の接続定義）
      |    |    |
      |    |    ├── /logger（共通ロガーの実装）
      |    |    |
      |    |    ├── /migrations（Bunのマイグレーション用SQL・スクリプト）
      |    |    |
      |    |    └── /observability（トレース取得用の設定など）
      |    |
      |    └── /presentation（プレゼンテーション層）
      |         |
      |         ├── /gen（OpenAPI定義から生成したGoコード）
      |         |
      |         ├── /handler（ハンドラー層）
      |         |
      |         └── /router（ルーター設定）
      |
      ├── /openapi（OpenAPIの定義）
      |    |
      |    ├── /components/schemas（コンポーネントのスキーマ定義）
      |    |
      |    ├── /paths（各種APIの定義）
      |    |
      |    └── openapi.yaml
      |
      └── /tests（インテグレーションテスト・e2eテスト用）
```
  
<br />
  
## ローカル開発環境構築

### 1. 環境変数ファイルをリネーム
  
```
cp ./.env.example ./.env
```  
  
### 2. コンテナのビルドと起動
  
```
docker compose build --no-cache
docker compose up -d
```  
  
### 3. マイグレーションの初期化
  
以下のコマンドを実行し、マイグレーション管理用テーブルを追加  
  
```
docker compose exec api go run cmd/migrate/main.go init
```
  
> ※もしテスト用DBに実行したい場合は、オプション「env ENV=testing」を追加し、次のようなコマンド「docker compose exec api env ENV=testing go run cmd/migrate/main.go init」を実行して下さい。
  
### 4. マイグレーション状態の確認
  
マイグレーション状態を確認したい場合は以下のコマンドを実行 
   
```
docker compose exec api go run cmd/migrate/main.go status
```
  
### 5. マイグレーションの実行
  
```
docker compose exec api go run cmd/migrate/main.go migrate
```
> ※マイグレーション実行後、一度に実行したSQLを一つのグループとして管理しています。
  
### 6. コンテナの停止・削除
  
```
docker compose down
```  
> ※ボリュームも合わせて削除したい場合は、オプション「-v」を付けて実行して下さい。（例：docker compose down -v）  
  
<br/>
  
## マイグレーション用のコマンド

### 1. マイグレーションファイルの新規作成
  
以下のコマンドを実行し、upとdown用の二つのファイルを作成します。  
upがマイグレーション実行用で、downがロールバック用です。  
  
```
docker compose exec api go run cmd/migrate/main.go create_sql [ファイル名]
```
> ※ファイル名の例「create_users_table」
  
### 2. マイグレーションのロールバック
```
docker compose exec api go run cmd/migrate/main.go rollback
```
  
> ※直前に実行したマイグレーションを一つのグループ単位とし、グループ単位で一つ前に戻します。  
  
<br/>
  
## oapi-codegenでGoのコードを生成するコマンド

OpenAPI定義は `src/openapi` 配下にあります。  
OpenAPI定義を修正した際は以下のコマンドを実行し、Goのコードを生成（上書き）して下さい。
  
```
make generate
```
  
<br/>
  
## mockgenでインターフェース定義からモック用のコードを生成するコマンド
  
リポジトリなどテストコードでモック化が必要な場合、
mockgenを利用してインターフェース定義からモック用のコード生成が可能です。
  
mockgenを利用したい場合、以下のコマンドを参考に
「-source」と「-destination」を指定して実行して下さい。
  
```
docker compose run --rm api mockgen -source=./internal/shared/logger/logger.go -destination=./internal/shared/logger/mock_logger/mock_logger.go
```
  
<br/>
  
## コード修正後に使うコマンド
  
### 1. go.modの修正
  
```
docker compose run --rm api go mod tidy
```  
  
### 2. フォーマット修正
  
```
docker compose run --rm api golangci-lint fmt -v ./...
```  
  
### 3. コード解析チェック
  
```
docker compose run --rm api golangci-lint run -v ./...
```  
  
### 4. テストコードの実行
  
テストコードのファイル（ _test.go ）を追加したパッケージのみテストを実行（ビルドタグ指定あり）  
  
・ユニットテスト実行  
  
```
make test-unit
```
  
・インテグレーションテスト実行  
  
```
make test-integration
```
  
・e2eテスト実行  
  
```
make test-e2e
```
  
・全てのテスト実行
  
```
make test
```
  
<br>
  
## 参考記事  
[]()  
  
作成中。。。