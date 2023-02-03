package policy

import "git.garena.com/shopee/loan-service/credit_backend/fast-escrow/go-retry/retry"

// AlwaysRetryPolicy 一种无穷重试策略
type AlwaysRetryPolicy struct {
	NeverRetryPolicy
}

func (a *AlwaysRetryPolicy) CanRetry(ctx retry.RtyContext) bool {
	return true
}
