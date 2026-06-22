# Member Point Calculation Module

Issue #1 の会員ポイント計算機能は `internal/core/member` に配置する。

`docs/rules/module-classification.md` の core 条件に照らし、会員ランク別の付与率、購入金額 10,000 円以上の倍率、小数点以下切り捨てという業務ルールを持つためである。

単純 CRUD や技術的処理ではなく、会員取得、入力検証、ポイント計算、エラー分類の複数ステップを扱うため generic には分類しない。単一手続きとして閉じる supporting より業務意味が強く、Issue でも DDD 実装先が `internal/core/member` と明示されている。
