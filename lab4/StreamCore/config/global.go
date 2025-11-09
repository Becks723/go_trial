package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *config
)

/* singleton */
func Instance() *config {
	once.Do(func() {
		instance = getInstance()
	})
	return instance
}

func getInstance() *config {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("error reading config: ", err.Error()) // TODO: i18n
	}

	var c config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatal("error resolve config: ", err.Error()) // TODO: i18n
	}
	return &c
}
