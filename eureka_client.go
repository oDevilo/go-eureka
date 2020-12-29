package eureka_client

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	STARTING = "starting"
	RUNNING  = "running"
	STOPING  = "stoping"
)

// eureka客户端
type EurekaDiscoveryClient struct {
	mutex   sync.RWMutex
	Running bool
	Config  *EurekaClientConfig
	// eureka服务中注册的应用
	Applications *Applications
	// 本服务信息
	Instance *Instance
}

// 启动时注册客户端，并后台刷新服务列表，以及心跳
func (c *EurekaDiscoveryClient) Start() error {
	c.mutex.Lock()
	if c.Running {
		return nil
	}
	c.Running = true
	// 注册
	if err := c.register(); err != nil {
		log.Println("eureka client start fail " + err.Error())
		return err
	}
	c.mutex.Unlock()
	log.Println("register application instance successful")
	// 刷新服务列表
	go c.refreshTask()
	// 心跳
	go c.heartbeatTask()
	// 监听退出信号，自动删除注册信息
	go c.handleSignal()
	return nil
}

// 刷新服务列表
func (c *EurekaDiscoveryClient) refreshTask() {
	for {
		if c.Running {
			if err := c.refresh(); err != nil {
				log.Println(err)
			} else {
				log.Println("refreshTask application instance successful")
			}
		} else {
			break
		}
		sleep := time.Duration(c.Config.RegistryFetchIntervalSeconds)
		time.Sleep(sleep * time.Second)
	}
}

// 心跳
func (c *EurekaDiscoveryClient) heartbeatTask() {
	for {
		if c.Running {
			if err := c.doHeartbeat(); err != nil {
				log.Println(err)
			} else {
				log.Println("heartbeat application instance successful")
			}
		} else {
			break
		}
		sleep := time.Duration(c.Config.RenewalIntervalInSecs)
		time.Sleep(sleep * time.Second)
	}
}

func (c *EurekaDiscoveryClient) register() error {
	return AddInstance(c.Config.ServiceUrl, c.Config.App, c.Instance)
}

func (c *EurekaDiscoveryClient) unRegister() error {
	instance := c.Instance
	return DeleteInstance(c.Config.ServiceUrl, instance.App, instance.InstanceID)
}

func (c *EurekaDiscoveryClient) doHeartbeat() error {
	instance := c.Instance
	return InstanceStatus(c.Config.ServiceUrl, instance.App, instance.InstanceID)
}

func (c *EurekaDiscoveryClient) refresh() error {
	applications, err := ListApplications(c.Config.ServiceUrl)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	c.Applications = applications
	c.mutex.Unlock()
	return nil
}

// 监听退出信号，删除注册的实例
func (c *EurekaDiscoveryClient) handleSignal() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		switch <-signalChan {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("receive exit signal, client instance going to de-register")
			err := c.unRegister()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("unRegister application instance successful")
			}
			os.Exit(0)
		}
	}
}

// 创建客户端
func NewClient(config *EurekaClientConfig) (*EurekaDiscoveryClient, error) {
	initConfig(config)
	ip, err := GetLocalIP()
	if err != nil {
		return nil, err
	}
	return &EurekaDiscoveryClient{Config: config, Instance: NewInstance(ip, config)}, nil
}

func initConfig(config *EurekaClientConfig) {
	if config.RenewalIntervalInSecs == 0 {
		config.RenewalIntervalInSecs = 30
	}
	if config.RegistryFetchIntervalSeconds == 0 {
		config.RegistryFetchIntervalSeconds = 15
	}
	if config.DurationInSecs == 0 {
		config.DurationInSecs = 90
	}
	config.App = strings.ToUpper(config.App)
	if config.Port == 0 {
		config.Port = 80
	}
}
