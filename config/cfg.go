package config

import (
	"github.com/lancer-kit/armory/tools"
)

type BaseConfig struct {
	Net            *Net           `yaml:"net"`
	LoadTestConfig LoadTestConfig `yaml:"load_tests"`
	NATS           string         `yaml:"nats"`
}

type Net struct {
	Timeout         int    `yaml:"timeout"`
	KeepAlive       int    `yaml:"keep_alive"`
	IdleConnTimeout int    `yaml:"idle_conn_timeout"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxConnsPerHost int    `yaml:"max_conns_per_host"`
	RealIP          string `yaml:"X-Real-IP"`
	ForwardedForIP  string `yaml:"X-Forwarded-For"`
}

func (Net) Default() *Net {
	return &Net{
		Timeout:         15,
		KeepAlive:       30,
		MaxIdleConns:    100,
		IdleConnTimeout: 90,
		MaxConnsPerHost: 1000,
		RealIP:          "",
		ForwardedForIP:  "",
	}
}

func (n *Net) Get() *Net {
	if n != nil {
		return n
	}

	return Net{}.Default()
}

type LoadTestConfig struct {
	Users   int `yaml:"users"`
	Threads int `yaml:"threads"`

	ScriptsPercentage map[string]int `yaml:"scripts_percentage"`
}

type Service struct {
	Url        tools.URL `yaml:"url"`
	PathPrefix string    `yaml:"path_prefix"`
}

func (s Service) URLWithPath(path string) string {
	s.Url.SetBasePath(s.PathPrefix)
	return s.Url.WithPath(path)
}
