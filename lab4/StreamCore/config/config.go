package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	Server  serverConfig  `mapstructure:"server"`
	MySql   mySqlConfig   `mapstructure:"mysql"`
	General generalConfig `mapstructure:"general"`
	Etcd    *etcd
	Service *svc
}

var (
	once     sync.Once
	instance *config
)

func Init(serviceName string) {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName("config")
		v.AddConfigPath("./config")

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("config: error viper.ReadInConfig: %v", err)
		}

		instance = new(config)
		if err := v.Unmarshal(&instance); err != nil {
			log.Fatalf("config: error viper.Unmarshal: %v", err)
		}
		instance.Service = getService(serviceName, v)
	})
}

func Instance() *config {
	return instance
}
func getService(service string, v *viper.Viper) *svc {
	s := new(svc)
	s.Name = v.GetString("services." + service + ".name")
	s.AddrList = v.GetStringSlice("services." + service + ".addr")

	switch service {
	// future: extension code goes here
	}
	return s
}

type mySqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type serverConfig struct {
	Port string `mapstructure:"port"`
}

type generalConfig struct {
	PageSize int `mapstructure:"page_size"`
}

type etcd struct {
	Addr string
}

type svc struct {
	Name     string
	AddrList []string
}

// extension: subtype svc and add more fields
// extension: add method for config to get specified svc type
// func (c *config) ServiceAsXXX() *XXXSvc {
//
// }
