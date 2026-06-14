# 依存関係ルール定義

このファイルは「モジュール間の依存関係（importルール）」を定義します。

設計思想ではなく、実装時の強制ルールとして扱うこと。

## 1. 基本原則

本プロジェクトは完全分離構造を採用します。

- `internal/core`
- `internal/supporting`
- `internal/generic`
- `internal/shared`

これらは明確に分離される。

## 2. 依存の基本ルール

### 許可される依存

- `internal/core` → `internal/shared` のみ
- `internal/supporting` → `internal/shared` のみ
- `internal/generic` → `internal/shared` のみ

### internal/shared の特性

- `internal/shared` は他に依存してはいけない（完全独立）
- `internal/shared` は外部パッケージへの依存を禁止する
  - 使用可能なのはGo標準ライブラリのみ
  - 例：context / errors / fmt / time など
- `internal/shared` は全モジュールから利用可能

## 3. モジュール間依存ルール

### internal/core

- 他のモジュール（supporting/generic）には依存しない
- `internal/shared` のみ利用可能

### internal/supporting

- 他のモジュール（core/generic）には依存しない
- `internal/shared` のみ利用可能

### internal/generic

- 他のモジュール（core/supporting）には依存しない
- `internal/shared` のみ利用可能

### internal/shared

- いかなるモジュールにも依存してはいけない
- 完全に独立した共通基盤である

## 4. 禁止事項（重要）

- `internal/core` / `internal/supporting` / `internal/generic` 間のimportは禁止
- `internal/shared` 以外への依存は禁止
- 循環依存は禁止
- ビジネスロジックをsharedに書くことは禁止
- infrastructure的実装を `internal/shared` に入れることは禁止

## 5. 依存構造の最終形

依存関係は以下の形のみ許可される：

- `internal/core` → `internal/shared`
- `internal/supporting` → `internal/shared`
- `internal/generic` → `internal/shared`
- `internal/shared` → どこにも依存しない

## 6. 設計意図

このルールの目的は以下である：

- モジュール間の結合を完全に排除する
- 依存関係の事故を防ぐ
- AI生成時の構造崩壊を防ぐ
- 予測可能なコード生成を保証する

## 7. 判断基準

迷った場合は以下を優先する：

- 依存を持たないことを優先する
- sharedへの依存のみ許可する
- 他モジュール参照が必要な場合は設計を見直す

## 8. 例

### OK

- `internal/core` → `internal/shared` のみimportして利用可能（shared配下すべて対象）
- `internal/supporting` → `internal/shared` のみimportして利用可能（shared配下すべて対象）
- `internal/generic` → `internal/shared` のみimportして利用可能（shared配下すべて対象）
- `internal/shared` → どこにも依存しない

### NG

- `internal/core` → `internal/supporting`
- `internal/supporting` → `internal/generic`
- `internal/generic` → `internal/core`
