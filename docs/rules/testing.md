# テストのルール定義

本ファイルはテストの分類・責務・モック方針を定義する。
TDD（Red → Green → Refactor）を前提とする。

## 1. 基本方針

- 必ずTDDで開発する（Red → Green → Refactor）
- 実装より先にテストを書く（Red）
- テストは仕様の定義として扱う
- 外部依存をテストに持ち込まない設計を優先する

## 2. テスト分類

### unit test

最小単位のロジック検証。

対象：
- domain
- usecase
- supporting service
- genericの純粋関数・単純処理

特徴：
- DBアクセスなし
- 外部APIなし
- 高速実行
- 純粋ロジック中心

### integration test

複数コンポーネントを結合したテスト。

対象：
- repository
- infrastructure
- DBアクセスを含む処理

特徴：
- ローカルDockerのDBを使用する
- 実DBでSQL・ORM（Bun）を検証する
- 外部APIは使用しない（必ずモック化）

### e2e test

API全体の動作確認。

対象：
- HTTP API全体
- handler → usecase / service → repository / DBアクセス処理

特徴：
- 実際のHTTPリクエストを送る
- DBはDocker環境を使用
- 外部APIは必ずモック化

※ただしe2eは必須ではなく、重要な業務フローのみ対象とする

## 3. モックルール

### 共通ルール

- 外部APIは必ずモック化する
- テストの安定性を優先する

### repositoryモック

- unit testではrepositoryは必ずmock化する
- usecaseテストではDBアクセスを排除する

### usecaseモック（重要）

- handlerテストではusecaseを必ずmock化する
- supportingのhandlerテストではserviceを必ずmock化する
- handlerはHTTPレイヤーのみに責務を限定する

## 4. テスト配置ルール

### unit test

- domain/usecase/service/generic処理と同階層、または近傍に配置
  - 例：
    - `xxx_usecase_test.go`
    - `xxx_service_test.go`

### integration test

- `tests/integration`

### e2e test

- `tests/e2e`

## 5. handlerテストの特別ルール

handlerテストでは以下を必須とする：

- usecase / serviceは必ずモック化する
- HTTPリクエスト・レスポンスのみを検証する
- domainロジックは一切テストしない

例：

- `src/internal/presentation/handler/..._handler_test.go`

## 6. テスト責務の境界

### domain
- 純粋ロジックの検証

### usecase
- 業務フローの検証（repositoryはmock）

### handler
- HTTP層の検証（usecase / serviceはmock）

### integration
- DB含む結合検証（外部APIはmock）

## 7. 判断ルール

迷った場合は以下を優先する：

- unit testを優先する
- 外部依存を増やさない
- mock化できるものはすべてmock化する

## 8. e2eの扱い（重要）

e2eテストは「必須ではない」。

以下の場合のみ追加する：

- 金銭系（決済など）
- 重要な業務フロー
- システム全体の整合性確認が必要な場合

それ以外はintegration testで十分とする
