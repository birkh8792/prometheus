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

func InitClient() naming_client.INamingClient {
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.18.27",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "7ec73a7e-ecba-49ef-a01a-65a288969ded", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "D:/tmp/nacos/log",
		CacheDir:            "D:/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
	return client
}

func TestRegisterService(t *testing.T) {
	client := InitClient()

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

func TestDeRegisterService(t *testing.T) {
	client := InitClient()

	param := vo.DeregisterInstanceParam{
		Ip:          "192.168.18.24",
		Port:        9100,
		ServiceName: "192.168.18.24:9100.json",
		Cluster:     "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
		Ephemeral:   true,            //it must be true
		//GroupName: "DEFAULT_GROUP",

	}
	success, _ := client.DeregisterInstance(param)
	fmt.Printf("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func TestGetService(t *testing.T) {
	client := InitClient()

	service, _ := client.GetService(vo.GetServiceParam{
		ServiceName: "192.168.18.24:9100.json",
		Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP",     // 默认值DEFAULT_GROUP
	})

	fmt.Printf("GetService, result:%+v \n\n", service)
}
