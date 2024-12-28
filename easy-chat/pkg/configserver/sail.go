package configserver

import (
	"encoding/json"
	"fmt"

	"github.com/HYY-yu/sail-client"
)

// Config 定义了配置服务器所需的配置参数
type Config struct {
	ETCDEndpoints  string `toml:"etcd_endpoints"`
	ProjectKey     string `toml:"project_key"`      // 项目密钥
	Namespace      string `toml:"namespace"`        // Etcd 命名空间
	Configs        string `toml:"configs"`          // 配置文件名
	ConfigFilePath string `toml:"config_file_path"` // 本地配置文件存放路径，空代表不存储本地配置文件
	LogLevel       string `toml:"log_level"`        // 日志级别(DEBUG\INFO\WARN\ERROR)，默认 WARN
}

type Sail struct {
	*sail.Sail
	sail.OnConfigChange
	c *Config
}

func NewSail(cfg *Config) *Sail {
	return &Sail{
		c: cfg,
	}
}

func (s *Sail) Build() error {
	var opts []sail.Option
	if s.OnConfigChange != nil {
		opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}
	// 创建并配置 Sail 实例
	s.Sail = sail.New(&sail.MetaConfig{
		ETCDEndpoints:  s.c.ETCDEndpoints,  // Etcd 端点
		ProjectKey:     s.c.ProjectKey,     // 项目密钥
		Namespace:      s.c.Namespace,      // Etcd 命名空间
		Configs:        s.c.Configs,        // 配置文件名
		ConfigFilePath: s.c.ConfigFilePath, // 配置文件路径（先删除再加载）
		LogLevel:       s.c.LogLevel,       // 日志级别
	}, opts...)
	return s.Sail.Err()
}

func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}
	return s.fromJsonBytes(s.Sail)
}

func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}

func (s *Sail) SetOnChange(f OnChange) {
	// 设置热加载方法
	s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
		data, err := s.fromJsonBytes(sail)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = f(data); err != nil {
			fmt.Println("OnChange err: ", err)
		}
	}
}
