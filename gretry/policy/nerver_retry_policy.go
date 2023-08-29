package policy

import (
	"github.com/liuhailove/go-retry/retry/context"
)

// NeverRetryPolicy 允许第一次尝试，不允许之后的重试
type NeverRetryPolicy struct {
}

func (n *NeverRetryPolicy) CanRetry(ctx gretry.RtyContext) bool {
	return (ctx.(*NeverRetryContext)).IsFinished()
}

func (n *NeverRetryPolicy) Open(parent gretry.RtyContext) gretry.RtyContext {
	var ntc = &NeverRetryContext{}
	ntc.Parent = parent
	return ntc
}

func (n *NeverRetryPolicy) Close(ctx gretry.RtyContext) {
	// no-op
}

func (n NeverRetryPolicy) RegisterError(ctx gretry.RtyContext, err error) {
	(ctx.(*NeverRetryContext)).setFinished()
	(ctx.(*NeverRetryContext)).RegisterError(err)
}

type NeverRetryContext struct {
	context.RtyContextSupport
	gretry.SimpleAttributeAccessorSupport
	Finished bool
}

func (n *NeverRetryContext) IsFinished() bool {
	return n.Finished
}

func (n *NeverRetryContext) setFinished() {
	n.Finished = true
}

func NewNeverRetryPolicy() *NeverRetryPolicy {
	return &NeverRetryPolicy{}
}
