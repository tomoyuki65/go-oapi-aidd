# TDD開発フロー

本フローは pm の指示を起点として進行するが、
各エージェントは自身の責務に従い独立して実行する。

本プロジェクトは「OpenAPI駆動 + TDD（Red → Green → Refactor）」を前提とする。

## 1. 開発フロー

1. pmが要件を整理し、タスクとして分解する  
2. testerが失敗するテストを書く（RED）  
3. reviewerがテスト内容をレビューする（REDレビュー）  
4. implementerがテストを通す（GREEN）  
5. reviewerが実装をレビューする（GREENレビュー）  
6. 問題があれば修正ループに戻る  

## 2. API開発時の追加フロー（重要）

API開発の場合は以下を必ず先に行う：

### ① OpenAPI定義の更新

- `src/openapi/openapi.yaml`
- `src/openapi/paths/*`
- `src/openapi/components/schemas/*`

ここでAPI仕様を確定させる（実装より優先）。

### ② コード生成

OpenAPI更新後は必ず生成する。

```bash
make generate
```

生成コード：

- `src/internal/presentation/gen`

※生成コードは編集禁止

## 3. テストの責務ルール

### testerの責務

- RED（失敗するテスト）を先に書く
- API仕様の抜けを検出する
- handler / usecase の境界を明確化する

### テストのモックルール

- handler → usecase は必ずmock化する
- usecase → repository は必ずmock化する
- 外部APIは必ずmock化する
- DBはintegration test以外では使用しない

## 4. 実装ルール（GREEN）

実装に関する詳細ルールは以下を参照する：

- `src/AGENTS.md`

## 5. レビュー（reviewer）

reviewerは以下を検査する：

### REDレビュー
- テストの妥当性
- 仕様の過不足
- モックの適切性
- API設計の整合性

### GREENレビュー
- 依存ルール違反
- レイヤー逸脱
- 実装の妥当性
- テストとの一致
- OpenAPIとの整合性

## 6. ルール

- REDなしで実装しない
- REDレビューなしでGREENに進まない
- GREEN未達でレビューしない
- GREENレビューなしで完了しない
- レビューNGなら再実装
- OpenAPI未定義のAPI実装は禁止
- 生成コードは編集禁止
- DBアクセスをunit testに持ち込まない
- 外部APIの実通信は禁止

## 7. 判断ルール

判断に迷った場合は以下を参照する：

- `docs/rules/module-classification.md`
- `docs/rules/dependency-rules.md`
- `docs/rules/architecture.md`

それでも不明な場合は、最も軽い構造を優先する：

generic → supporting → core
