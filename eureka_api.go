package eureka_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://blog.csdn.net/qq_33475202/article/details/83654446
// eureka       /
// spring cloud /eureka/

type InstanceInfo struct {
	Instance *Instance `json:"instance"`
}

// 注册实例 POST /eureka/v2/apps/appID
func AddInstance(url, app string, instance *Instance) error {
	var info = &InstanceInfo{
		Instance: instance,
	}
	infoJson, err := json.Marshal(info)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, UrlAppend(url, "/apps/"+app), bytes.NewReader(infoJson))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("add application instance failed, error: %s", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d %s", resp.StatusCode, string(body))
	}
	return nil
}

// 删除实例 DELETE /eureka/v2/apps/appID/instanceID
func DeleteInstance(url, app, instanceID string) error {
	req, err := http.NewRequest(http.MethodDelete, UrlAppend(url, "/apps/"+app+"/"+instanceID), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete application instance failed, error: %s", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d %s", resp.StatusCode, string(body))
	}
	return nil
}

type Result struct {
	Applications *Applications `json:"applications"`
}

// 查询所有服务实例 GET /eureka/v2/apps
func ListApplications(url string) (*Applications, error) {
	apps := new(Applications)
	res := &Result{
		Applications: apps,
	}

	req, err := http.NewRequest(http.MethodGet, UrlAppend(url, "/apps"), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", " application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get applications failed, error: %s", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get applications failed, error: %s", err)
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, fmt.Errorf("get applications failed, error: %s", err)
	}
	return apps, nil
}

// 发送心跳 PUT /eureka/v2/apps/appID/instanceID
func InstanceStatus(eurekaUrl, app, instanceID string) error {
	Url, _ := url.Parse(UrlAppend(eurekaUrl, "/apps/"+app+"/"+instanceID))
	params := url.Values{
		"status": {"UP"},
	}
	Url.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodPut, Url.String(), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d %s", resp.StatusCode, string(body))
	}
	return nil
}
