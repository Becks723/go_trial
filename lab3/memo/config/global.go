package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var instance *config
var once sync.Once

/* Config单例 */
func Instance() *config {
	// once.Do保障线程安全
	once.Do(func() {
		instance = initConfig()
	})
	return instance
}

func initConfig() *config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("/app/config")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("err reading config") // TODO: i18n
	}

	var c *config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatal("err resolve config") // TODO: i18n
	}
	return c
}
