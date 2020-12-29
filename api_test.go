package eureka_client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	url := "http://test.local.com:8001/"
	applications, err := ListApplications(url)
	if err != nil {
		t.Log(err)
	}
	marshal, _ := json.Marshal(applications)
	t.Log(string(marshal))

	InstanceStatus("http://test.local.com:8001/", "kps-gateway", "10.242.28.186:kps-gateway:8080")
}

func TestRequest(t *testing.T) {
	var ApplicationCounter map[string]int
	if service, ok := ApplicationCounter["a"]; ok {
		t.Log(service)
	} else {
		t.Log("has not")
	}

	ApplicationCounter["a"] = 1
	if service, ok := ApplicationCounter["a"]; ok {
		t.Log(service)
	} else {
		t.Log("has not")
	}
}

func TestClient(t *testing.T) {
	client, _ := NewClient(&EurekaClientConfig{
		ServiceUrl: " http://10.219.192.172:8001/eureka/",
		App:        "cmdb",
		Port:       3333,
	})
	if err := client.Start(); err != nil {
		t.Log(err)
		return
	}

	time.Sleep(time.Second)

	httpRequestClient := HttpRequestClient{
		RetriesTime:  2,
		Interceptors: &[]HttpRequestInterceptor{NewEurekaHttpRequestInterceptor(client)},
	}
	req, _ := http.NewRequest(http.MethodGet, "http://doorkeeper-server/login/ping", nil)
	resp, err := httpRequestClient.Do(req)
	if err != nil {
		t.Log(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	t.Log(string(body))
}
