# oapi-codegenに関するルール定義

このファイルはOpenAPI定義およびoapi-codegenに関するルールを定義する。

## 1. 基本方針

本プロジェクトではAPI実装より先にOpenAPIを定義する。

API仕様を唯一の正とし、GoコードはOpenAPIから生成する。

実装からAPI仕様を作成してはならない。

## 2. OpenAPI First

新規API追加・既存API変更時は以下の順序を厳守する。

1. OpenAPI定義を修正する
2. コード生成を実行する
3. テストコード（Red）を作成する
4. 実装する（Green）
5. リファクタリングする（Refactor）

## 3. OpenAPI定義の配置

API定義は以下に配置する。

- `src/openapi/openapi.yaml`
- `src/openapi/paths/*`
- `src/openapi/components/schemas/*`

API仕様は必ずOpenAPI上で管理する。

## 4. コード生成

OpenAPI定義変更後は必ず以下を実行する。

```bash
make generate
```

生成コードは以下に出力される。

- `src/internal/presentation/gen`

## 5. 生成コードの扱い

生成コードは編集してはならない。

禁止事項：

- hand edit
- コメント追加
- メソッド追加
- 構造体修正
- import追加

変更が必要な場合は OpenAPI 定義を修正し再生成する。

## 6. 使用範囲

生成コードは presentation 層でのみ利用する。

利用可能：

- `src/internal/presentation/handler`
- `src/internal/presentation/router`

利用不可：

- `src/internal/core`
- `src/internal/supporting`
- `src/internal/generic`
- `src/internal/shared`

## 7. DTOルール

生成される request / response は DTO として扱う。

DTOは外部境界専用であり、業務モデルとして利用してはならない。

## 8. ドメインモデルとの分離

DTOと業務モデルは分離する。

禁止事項：

- DTOを永続化モデルとして利用する
- DTOを業務モデルとして利用する
- DTOをそのまま内部処理へ渡す

handler層で変換を行うこと。

## 9. handlerの責務

handlerは以下のみを担当する。

- リクエスト受信
- DTO変換
- usecase / service呼び出し
- レスポンス変換

業務ロジックを書いてはならない。

## 10. 判断基準

API変更が必要な場合は必ず以下を確認する。

- OpenAPIを先に変更したか
- make generate を実行したか
- 生成コードを編集していないか

迷った場合はOpenAPI定義を正として扱う。
