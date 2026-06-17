# アーキテクチャ設計のルール定義

このファイルは本プロジェクトにおける「構造の意味」を定義します。

コードの配置判断や依存ルールは扱いません。それらはそれぞれ以下に委譲：

- `docs/rules/module-classification.md`（配置判断）
- `docs/rules/dependency-rules.md`（依存関係）

## 1. 本プロジェクトのディレクトリ構成

```
/my-project
 |
 └── /src
      |
      ├── /cmd
      |    |
      |    ├── /migrate/main.go（DBのマイグレーション用スクリプト）
      |    |
      |    └── /seed/main.go（マスタデータ登録用スクリプト）
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
      |    ├── /generic（一般的な業務領域）
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
      |    |    |    ├── /seed（マスタデータ登録用のseed定義）
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

## 2. アーキテクチャの基本思想

本プロジェクトは以下を前提として設計される：

- OpenAPI定義からのGoコード生成（oapi-codegen）
- TDD（Test Driven Development）の徹底
- ドメイン志向のモジュール分割
- 依存関係を厳密に制御した制約ベースアーキテクチャ

本プロジェクトは「DDDの概念を一部取り入れた構造」であるが、
DDDの厳密なレイヤー構造・依存ルールには従わない。

## 3. 設計の位置づけ

本アーキテクチャは以下の中間的性質を持つ：

- DDDのように業務中心の構造を持つ
- ただし依存制御はDDDではなくルールベース
- Clean Architectureのようなレイヤー依存は持たない
- AI生成に最適化された単純構造を採用する

## 4. モジュールの役割定義

### core（中核の業務領域）

システムのビジネス価値そのもの。

- 業務ルールを持つ
- 状態変化を扱う
- 複数ステップの整合性を保証する
- 最も重要な制約を持つ

---

### supporting（補完的な業務領域）

単一業務フローを処理する領域。

- トランザクションスクリプト型
- 手続き的処理
- 軽量な業務ロジック
- ドメイン分割不要

---

### generic（一般的な業務領域）

業務意味を持たない処理。

- CRUD
- データ変換
- 文字列・時間・ID処理
- 純粋な技術関数

---

### shared（横断関心・共通基盤）

全モジュール共通の基盤。

- ログ
- エラー定義
- DTO
- ユーティリティ
- sharedは全モジュールから利用可能な唯一の共通基盤
- sharedは「業務ロジックを持たない純粋な共通部品のみ」を配置する
- sharedは「依存の終端」であり、上位概念を持たない

---

## 5. 設計思想

本アーキテクチャの目的は以下である：

- 責務の明確な分離
- 判断の機械化
- AI生成の安定化
- 構造の単純化

## 6. レイヤーの考え方

本プロジェクトでは「レイヤー」ではなく「責務領域」として扱う。

- core = 業務の中心
- supporting = 業務フロー
- generic = 技術処理
- shared = 共通基盤

機能の重要度の上下関係はあるが、依存関係の上下構造は存在しない。

## 7. 非目標

本設計は以下を目的としない：

- 厳密なDDDの再現
- Clean Architectureの完全準拠
- 柔軟性の最大化

## 8. 判断の委譲

詳細な判断は以下に委譲する：

- 配置 → `docs/rules/module-classification.md`
- 依存 → `docs/rules/dependency-rules.md`

## 9. 判断基準（概要）

- 業務の中心 → core
- 業務フロー → supporting
- 技術処理 → generic
- 横断基盤 → shared
