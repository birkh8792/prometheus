package nacos

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

var (
	// DefaultSDConfig is the default Consul SD configuration.
	DefaultSDConfig = SDConfig{
		Server: "localhost:8848",
	}
)

func init() {
	fmt.Println("nacos init...")
	discovery.RegisterConfig(&SDConfig{})
}

// SDConfig Name()  NewDiscoverer(opts discovery.DiscovererOptions) (discovery.Discoverer, error)
type SDConfig struct {
	Server              string `yaml:"server,omitempty"`
	ContextPath         string `yaml:"contextPath,omitempty"`
	Port                uint64 `yaml:"port,omitempty"`
	Scheme              string `yaml:"scheme,omitempty"`
	NamespaceId         string `yaml:"namespaceId,omitempty"`
	TimeoutMs           uint64 `yaml:"timeoutMs,omitempty"`
	NotLoadCacheAtStart bool   `yaml:"notLoadCacheAtStart,omitempty"`
	LogDir              string `yaml:"logDir,omitempty"`
	CacheDir            string `yaml:"cacheDir,omitempty"`
	RotateTime          string `yaml:"rotateTime,omitempty"`
	MaxAge              int64  `yaml:"maxAge,omitempty"`
	LogLevel            string `yaml:"logLevel,omitempty"`
	Username            string `yaml:"username,omitempty"`
	Password            string `yaml:"password,omitempty"`
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
	configClient *config_client.IConfigClient

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
	//get config
	//content, err := configClient.GetConfig(vo.ConfigParam{
	//	DataId: "chixiao-test.json",
	//	Group:  "DEFAULT_GROUP",
	//})
	//fmt.Println("GetConfig,config :" + content)

	if err != nil {
		return nil, err
	}

	cd := &Discovery{
		configClient: &configClient,
		logger:       logger,
	}

	return cd, nil
}

// Run implements the Discoverer interface.
func (d *Discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	fmt.Println("run invoke...")

}
