# Go-Retry: 高可用的Retry重试逻辑和退避的 Go 库。

它具有高度可扩展性，可以完全控制重试的方式和时间。您还可以通过实现 Backoff 接口来编写自己的自定义退避函数。

# 特点

- 可扩展：你可以自定义重试及会退逻辑。此方案的灵感来自于Spring-Retry
- 无依赖：除 Go 标准库依赖外，无其他依赖
- 高并发：所有逻辑都是并发安全的
- 模板化：提供了重试模板，可以快速定制重试会退策略
- 策略丰富：提供了多种重试及会退策略

# 策略讲解

## 重试策略：定义在被调方出现错误时怎么进行重试
重试策略：

- 1、NeverRetryPolicy：从不重试，这种策略无论在业务是否发生异常时，都不进行重试

- 2、SimpleRetryPolicy：简单重试，这种策略可以配置重试的次数，只要发生符合逻辑的错误，就会进行立即重试，而不进行重试规避

- 3、MaxAttemptsRetryPolicy：最大重试次数策略，这种策略若业务发生异常，只要重试次数小于配置的策略则进行重试

- 4、TimeoutRtyPolicy：超时重试策略，若业务接口出现错误，且接口没有没有超过对应的超时时间则重试

- 5、AlwaysRetryPolicy：无穷重试，只要接口出现错误，则进行重试

- 6、ErrorClassifierRetryPolicy：异常分类重试，错误符合对应的异常策略则进行重试

- 7、CompositeRetryPolicy：组合重试策略，多种重试策略的组合

## 回退策略：定义在错误发生后，下次重试前应该怎么处理

- 1、NoBackOffPolicy：无任何回退，直接重试

- 2、FixedBackOffPolicy：固定时间间隔回退策略

- 3、ExponentialBackOffPolicy：指数回退策略，需要设置最小回退值，递增倍率，最大回退值

- 4、ExponentialRandomBackOffPolicy：具备随机倍率的指数回退策略，和ExponentialBackOffPolicy相比，每次递增的倍率为[1，递增倍率]的随机整数

- 5、UniformRandomBackoffPolicy：均匀回退策略，需要配置最小回退值，最大回退值，每次回退时是[最小回退值，最大回退值]的随机值

## 用于重试/不重试的异常：
如果两者都为空，则表示任意异常都进行重试，两个配置不能同时满足

## 异常匹配：默认任意异常都会进行重试，
目前支持如下：

- 1、ExactMatch：精确匹配

- 2、PrefixMatch：前缀匹配

- 3、SuffixMatch：后缀匹配

- 4、ContainMatch：包含匹配

- 5、RegularMatch：正则匹配

- 6、AnyMatch：只要不为空，则匹配

# 应用场景

- 场景1：接口显示的抛出timeout异常

- 场景2：接口正常返回，但是返回值包含错误码

# 注意事项

以下场景建议不要设置自动重试规则：

- 1、不要对致命错误（Fatal Error）重试，重试条件策略一般采取白名单机制（如只对超时异常、限流异常重试）。 典型的报错像Error、Interface not exist这些都建议您不要进行重试，业务错误建议您不要进行重试。

- 2、重试需要注意幂等性，非幂等的操作在部分场景下重试（如超时，请求已到达对端）可能会不符合预期。

- 3、前端应用，超时时间长且慢的接口，建议您不要进行重试。

# 案例分析

``` Golang
func TestRetryTemplateBuilder(t *testing.T) {  
   var retryTemplate = NewRetryTemplateBuilder().  
       MaxAttemptsRtyPolicy(5).  
       //FixedBackoff(1000).  
       //ExponentialBackoffWithRandom(1000, 2, 5000, true).  
       //WithinMillisRtyPolicy(1000).  
       //InfiniteRtyPolicy().  
       UniformRandomBackoff(100, 1000).  
       //NotRetryOn(errors.New("error")).  
       RetryOn(errors.New("hello world")).  
       RetryOn(errors.New("error")).  
       WithErrorMatchPattern(classify.RegularMatch).  
       Build()  

   var result, err = retryTemplate.Execute(&MyTestRetryCallback{})
   if err == nil {
       fmt.Println(result)
   } else {
       fmt.Println(err)
   }

}

type MyTestRetryCallback struct {  
}

func (m MyTestRetryCallback) DoWithRetry(content retry.RtyContext) interface{} {  
	fmt.Println(content.GetRetryCount())  
	var result, err = PrintHello()  
	if err != nil {  
		panic(err.Error())  
	}  
	return result  
}

var i = 0

func PrintHello() (string, error) {  
	if i < 2 {  
		i++  
		if i == 1 {  
			panic("error")  
		}  
		util.Sleep(time.Millisecond * 100)  
		fmt.Println("PrintHello error")  
		return "", errors.New("error")  
	} else {  
		//fmt.Println("hello")  
		return "hello world", nil  
	}  
}   
```