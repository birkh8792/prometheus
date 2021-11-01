package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
)

func TestConfiguredService(t *testing.T) {
	fmt.Println("test begin...")
}

func initClient() naming_client.INamingClient {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "7ec73a7e-ecba-49ef-a01a-65a288969ded", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "D:/tmp/nacos/log",
		CacheDir:            "D:/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "192.168.18.27",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}
	return namingClient
}

func TestRegisterService(t *testing.T) {
	client := initClient()

	param := vo.RegisterInstanceParam{
		Ip:          "192.168.18.24",
		Port:        9100,
		ServiceName: "192.168.18.24:9100.json",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
		Metadata: map[string]string{
			"__meta_chixiao_hostName":            "",
			"__meta_chixiao_operatingSystem":     "",
			"__meta_chixiao_team":                "",
			"__meta_chixiao_hostLabel":           "",
			"__meta_chixiao_project":             "",
			"__meta_chixiao_port":                "9100",
			"__meta_chixiao_hostLabelName":       "",
			"__meta_chixiao_cxInstance":          "192.168.18.24:9100",
			"__meta_chixiao_app":                 "chixiao",
			"__meta_chixiao_ip":                  "192.168.18.24",
			"__meta_chixiao_businessType":        "MYSQL",
			"__meta_chixiao_operatingSystemName": "",
			"__meta_chixiao_serveName":           "",
			"__meta_chixiao_hostGroupName":       "",
			"__meta_chixiao_agentPort":           "9100",
			"__meta_chixiao_domainName":          "云HIS",
			"__meta_chixiao_hostGroup":           "",
			"__meta_chixiao_tags":                "database-export",
			"__meta_chixiao_domainId":            "1385394557942902723",
		},
	}

	success, _ := client.RegisterInstance(param)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func TestRegisterInstance(t *testing.T) {
	namingClient := initClient()
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.18.24",
		Port:        9100,
		ServiceName: "hostservice",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(success)
}

func TestDeregisterInstance(t *testing.T) {
	success, err := initClient().DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        9182,
		ServiceName: "127.0.0.1:9182",
		Ephemeral:   true,
		Cluster:     "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(success)
}
