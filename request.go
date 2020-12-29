package eureka_client

import (
	"net/http"
)

// 请求拦截器
type HttpRequestInterceptor interface {
	intercept(request *http.Request) (*http.Request, error)
}

type HttpRequestClient struct {
	// 请求拦截器 用于请求前处理
	Interceptors *[]HttpRequestInterceptor
	// 重试次数
	RetriesTime int
}

func (c *HttpRequestClient) Do(request *http.Request) (*http.Response, error) {
	interceptors := *(c.Interceptors)
	if len(interceptors) > 0 {
		for i := range interceptors {
			if intercept, err := interceptors[i].intercept(request); err != nil {
				return nil, err
			} else {
				request = intercept
			}
		}
	}
	retriesTime := c.RetriesTime
	if retriesTime <= 0 {
		retriesTime = 1
	}
	var response *http.Response
	var err error
	for i := 0; i < retriesTime; i++ {
		response, err = http.DefaultClient.Do(request)
		if err == nil {
			return response, err
		}
	}
	return nil, err
}
