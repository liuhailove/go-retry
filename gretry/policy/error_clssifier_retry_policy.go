package policy

import (
	"github.com/liuhailove/go-retry/retry/classify"
	"github.com/liuhailove/go-retry/retry/context"
)

type ErrorClassifierRetryPolicy struct {
	ErrorClassifier *classify.ErrorClassifier
}

func (e ErrorClassifierRetryPolicy) CanRetry(ctx gretry.RtyContext) bool {
	var err = ctx.GetLastError()
	return err == nil || e.ErrorClassifier.Classify(err)
}

func (e ErrorClassifierRetryPolicy) Open(parent gretry.RtyContext) gretry.RtyContext {
	return &context.RtyContextSupport{Parent: parent}
}

func (e ErrorClassifierRetryPolicy) Close(ctx gretry.RtyContext) {
	// no-op
}

func (e ErrorClassifierRetryPolicy) RegisterError(ctx gretry.RtyContext, err error) {
	var simpleContext = ctx.(*context.RtyContextSupport)
	simpleContext.RegisterError(err)
}
