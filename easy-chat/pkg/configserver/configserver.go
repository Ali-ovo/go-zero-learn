package configserver

import (
	"errors"

	"github.com/zeromicro/go-zero/core/conf"
)

type OnChange func([]byte) error

var ErrNotSetConfig = errors.New("config not set")

type ConfigServer interface {
	Build() error
	SetOnChange(OnChange)
	FromJsonBytes() ([]byte, error)
}

type configServer struct {
	ConfigServer
	configFile string
}

func NewConfigServer(configFile string, s ConfigServer) *configServer {
	return &configServer{
		ConfigServer: s,
		configFile:   configFile,
	}
}

func (s *configServer) MustLoad(v any, onChange OnChange) error {
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}

	if s.ConfigServer == nil {
		conf.MustLoad(s.configFile, v)
		return nil
	}

	if onChange != nil {
		s.ConfigServer.SetOnChange(onChange)
	}

	if err := s.ConfigServer.Build(); err != nil {
		return err
	}

	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}

	// 从JSON字节加载配置
	return conf.LoadFromJsonBytes(data, v)
}

func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
