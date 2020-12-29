package eureka_client

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type EurekaHttpRequestInterceptor struct {
	EurekaDiscoveryClient *EurekaDiscoveryClient
	InstanceCounter       map[string]int
}

func (hr *EurekaHttpRequestInterceptor) intercept(request *http.Request) (*http.Request, error) {
	serviceName := request.URL.Host
	// 获取其中一个服务的真实ip:port
	applications := *(hr.EurekaDiscoveryClient.Applications.Applications)
	var application *Application
	for i := range applications {
		if strings.ToUpper(serviceName) == strings.ToUpper(applications[i].Name) {
			application = &applications[i]
			break
		}
	}
	if application == nil {
		return nil, errors.New("application is not found")
	}
	instances := *(application.Instances)
	if len(instances) <= 0 {
		return nil, errors.New("application instance is not found")
	}
	// 从实例里面轮询
	var minInstance Instance
	var minCount int
	for i := range instances {
		instance := instances[i]
		// 实例不是存活的，则下一个
		if "UP" != instance.Status {
			continue
		}
		if count, ok := hr.InstanceCounter[instance.InstanceID]; ok {
			if minCount == 0 || minCount > count {
				minInstance = instance
				hr.InstanceCounter[instance.InstanceID] = count + 1
			}
		} else {
			hr.InstanceCounter[instance.InstanceID] = 1
			minInstance = instance
			break
		}
	}
	// 利用这个实例 修改请求
	request.URL.Host = minInstance.HostName + ":" + strconv.Itoa(minInstance.Port.Port)
	return request, nil
}

func NewEurekaHttpRequestInterceptor(client *EurekaDiscoveryClient) *EurekaHttpRequestInterceptor {
	instanceCounter := make(map[string]int)
	return &EurekaHttpRequestInterceptor{
		EurekaDiscoveryClient: client,
		InstanceCounter:       instanceCounter,
	}
}
