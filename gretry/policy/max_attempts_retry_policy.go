package policy

import (
	"github.com/liuhailove/go-retry/retry/context"
)

const (
	DefaultMaxRetryAttempts int32 = 3
)

type MaxAttemptsRetryPolicy struct {
	// 最大重试次数
	MaxAttempts int32
}

func (m *MaxAttemptsRetryPolicy) CanRetry(ctx gretry.RtyContext) bool {
	return ctx.GetRetryCount() < m.MaxAttempts
}

func (m *MaxAttemptsRetryPolicy) Open(parent gretry.RtyContext) gretry.RtyContext {
	return &context.RtyContextSupport{Parent: parent}
}

func (m *MaxAttemptsRetryPolicy) Close(ctx gretry.RtyContext) {
	// no-op
}

func (m *MaxAttemptsRetryPolicy) RegisterError(ctx gretry.RtyContext, err error) {
	(ctx.(*context.RtyContextSupport)).RegisterError(err)
}

func NewMaxAttemptsRetryPolicy() *MaxAttemptsRetryPolicy {
	return &MaxAttemptsRetryPolicy{MaxAttempts: DefaultMaxRetryAttempts}
}

func NewMaxAttemptsRetryPolicyWithAttempts(maxAttempts int32) *MaxAttemptsRetryPolicy {
	return &MaxAttemptsRetryPolicy{MaxAttempts: maxAttempts}
}
