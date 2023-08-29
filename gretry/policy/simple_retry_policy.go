package policy

import (
	"errors"
	"github.com/liuhailove/go-retry/retry/classify"
	"github.com/liuhailove/go-retry/retry/context"
	"strconv"
)

var (
	// DefaultMaxAttempts 默认重试次数
	DefaultMaxAttempts int32 = 3
)

// SimpleRetryPolicy 根据具体的错误重试固定次数.
// 例如
//
//	retryTemplate = new RetryTemplate(new SimpleRetryPolicy(3));
//	retryTemplate.execute(callback);
//	上面的例子会至少重试一次，最多重试3次
type SimpleRetryPolicy struct {
	MaxAttempts         int32
	RetryableClassifier *classify.ErrorClassifier
}

func NewSimpleRetryPolicy() *SimpleRetryPolicy {
	return NewSimpleRetryPolicyWithMaxAttemptsAndErrors(DefaultMaxAttempts, []error{errors.New("any match")})
}

func NewSimpleRetryPolicyWithMaxAttemptsAndErrors(maxAttempts int32, errs []error) *SimpleRetryPolicy {
	return NewSimpleRetryPolicyWithMaxAttemptsAndErrorsAndDefault(maxAttempts, errs, false)
}

func NewSimpleRetryPolicyWithMaxAttemptsAndErrorsAndDefault(maxAttempts int32, errs []error, defaultValue bool) *SimpleRetryPolicy {
	var inst = new(SimpleRetryPolicy)
	inst.MaxAttempts = maxAttempts
	inst.RetryableClassifier = new(classify.ErrorClassifier)
	inst.RetryableClassifier.SetClassified(errs)
	inst.RetryableClassifier.DefaultValue = defaultValue
	return inst
}

func (s *SimpleRetryPolicy) CanRetry(ctx gretry.RtyContext) bool {
	var err = ctx.GetLastError()
	var meet bool
	if err == nil {
		meet = true
	} else {
		meet = s.RetryableClassifier.Classify(err)
	}
	return meet && ctx.GetRetryCount() < s.MaxAttempts
}

func (s *SimpleRetryPolicy) Open(parent gretry.RtyContext) gretry.RtyContext {
	return NewSimpleRetryContext(parent)
}

func (s *SimpleRetryPolicy) Close(ctx gretry.RtyContext) {
	// no-op
}

func (s *SimpleRetryPolicy) RegisterError(ctx gretry.RtyContext, err error) {
	var simpleContext = ctx.(*SimpleRetryContext)
	simpleContext.RegisterError(err)
}

func (s *SimpleRetryPolicy) String() string {
	return "SimpleRetryPolicy[maxAttempts=" + strconv.FormatInt(int64(s.MaxAttempts), 10) + "]"
}

type SimpleRetryContext struct {
	context.RtyContextSupport
	gretry.SimpleAttributeAccessorSupport
}

func NewSimpleRetryContext(parent gretry.RtyContext) *SimpleRetryContext {
	var inst = new(SimpleRetryContext)
	inst.Parent = parent
	return inst
}
