package nacos

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/discovery/refresh"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

var (
	// DefaultSDConfig is the default Nacos SD configuration.
	DefaultSDConfig = SDConfig{
		Server: "localhost:8848",
	}
)

func init() {
	discovery.RegisterConfig(&SDConfig{})
}

// SDConfig Name()  NewDiscoverer(opts discovery.DiscovererOptions) (discovery.Discoverer, error)
type SDConfig struct {
	Server              string         `yaml:"server,omitempty"`
	ContextPath         string         `yaml:"contextPath,omitempty"`
	Port                uint64         `yaml:"port,omitempty"`
	Scheme              string         `yaml:"scheme,omitempty"`
	NamespaceId         string         `yaml:"namespaceId,omitempty"`
	TimeoutMs           uint64         `yaml:"timeoutMs,omitempty"`
	NotLoadCacheAtStart bool           `yaml:"notLoadCacheAtStart,omitempty"`
	LogDir              string         `yaml:"logDir,omitempty"`
	CacheDir            string         `yaml:"cacheDir,omitempty"`
	RotateTime          string         `yaml:"rotateTime,omitempty"`
	MaxAge              int64          `yaml:"maxAge,omitempty"`
	LogLevel            string         `yaml:"logLevel,omitempty"`
	Username            string         `yaml:"username,omitempty"`
	Password            string         `yaml:"password,omitempty"`
	DataId              string         `yaml:"dataId,omitempty"`
	Group               string         `yaml:"group,omitempty"`
	RefreshInterval     model.Duration `yaml:"refresh_interval,omitempty"`
}

// Name returns the name of the Config.
func (*SDConfig) Name() string { return "nacos" }

// NewDiscoverer returns a Discoverer for the Config.
func (c *SDConfig) NewDiscoverer(opts discovery.DiscovererOptions) (discovery.Discoverer, error) {
	return NewDiscovery(c, opts.Logger)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if strings.TrimSpace(c.Server) == "" {
		return errors.New("nacos SD configuration requires a server address")
	}
	if c.Username == "" || c.Password == "" {
		return errors.New("at most one of consul SD configuration username and password and basic auth can be configured")
	}

	return nil
}

// Discovery implements the Discoverer interface for discovering
// targets from nacos.
type Discovery struct {
	*refresh.Discovery

	dataId string
	group  string

	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient

	sources map[string]*targetgroup.Group

	parse  func(data []byte, path string) (model.LabelSet, error)
	logger log.Logger
}

// NewDiscovery returns a new Discovery for the given config.
func NewDiscovery(conf *SDConfig, logger log.Logger) (*Discovery, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           conf.TimeoutMs,
		NotLoadCacheAtStart: conf.NotLoadCacheAtStart,
		LogDir:              conf.LogDir,
		CacheDir:            conf.CacheDir,
		RotateTime:          conf.RotateTime,
		MaxAge:              conf.MaxAge,
		LogLevel:            conf.LogLevel,
	}

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      conf.Server,
			ContextPath: conf.ContextPath,
			Port:        conf.Port,
			Scheme:      conf.Scheme,
		},
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	//get config
	//content, err := configClient.GetConfig(vo.ConfigParam{
	//	DataId: conf.DataId,
	//	Group:  conf.Group,
	//})
	//fmt.Println("GetConfig,config :" + content)

	if err != nil {
		return nil, err
	}

	cd := &Discovery{
		dataId:       conf.DataId,
		group:        conf.Group,
		configClient: configClient,
		namingClient: namingClient,
		logger:       logger,
	}

	cd.Discovery = refresh.NewDiscovery(
		logger,
		"nacos",
		time.Duration(conf.RefreshInterval),
		cd.refresh,
	)
	return cd, nil
}

func (d *Discovery) refresh(ctx context.Context) ([]*targetgroup.Group, error) {
	// 从nacos获取服务名列表:GetAllServicesInfo
	serviceInfos, err := d.namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "7ec73a7e-ecba-49ef-a01a-65a288969ded",
		PageNo:    1,
		PageSize:  10,
	})
	if err != nil {
		return nil, err
	}

	tg := &targetgroup.Group{
		Source: "nacos",
	}

	for _, dom := range serviceInfos.Doms {
		targets := targetsForDom(&dom, d.namingClient)
		tg.Targets = append(tg.Targets, targets...)
	}

	return []*targetgroup.Group{tg}, nil
}

func fetchServices(namingClient naming_client.INamingClient) {

}

func targetsForDom(dom *string, namingClient naming_client.INamingClient) []model.LabelSet {
	fmt.Println(*dom)

	//获取服务信息：GetService
	services, err := namingClient.GetService(vo.GetServiceParam{
		ServiceName: *dom,
		Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP",     // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetService, result:%+v \n\n", services)

	targets := make([]model.LabelSet, 0, len(services.Hosts))

	var targetAddress string
	host := services.Hosts[0]

	targetAddress = host.Ip + ":" + strconv.FormatUint(host.Port, 10)
	target := model.LabelSet{
		model.AddressLabel:  lv(targetAddress),
		model.InstanceLabel: lv(services.Hosts[0].Ip + ":" + strconv.FormatUint(services.Hosts[0].Port, 10)),
	}
	targets = append(targets, target)

	return targets
}

func lv(s string) model.LabelValue {
	return model.LabelValue(s)
}

// Run implements the Discoverer interface.
//func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
//	fmt.Println("run invoke...")
//
//	//client := d.configClient
//
//	//get config
//	content, err := d.configClient.GetConfig(vo.ConfigParam{
//		DataId: "chixiao-test.json",
//		Group:  "DEFAULT_GROUP",
//	})
//	fmt.Println("GetConfig,config :" + content)
//	if err != nil {
//
//	}
//	var targetGroups []*targetgroup.Group
//
//	if err := json.Unmarshal([]byte(content), &targetGroups); err != nil {
//
//	}
//}
