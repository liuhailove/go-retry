package policy

// AlwaysRetryPolicy 一种无穷重试策略
type AlwaysRetryPolicy struct {
	NeverRetryPolicy
}

func (a *AlwaysRetryPolicy) CanRetry(ctx gretry.RtyContext) bool {
	return true
}
