package env

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type env struct {
	DB_Host                  string  `mapstructure:"db_host"`
	DB_Port                  string  `mapstructure:"db_port"`
	DB_Username              string  `mapstructure:"db_username"`
	DB_Password              string  `mapstructure:"db_password"`
	DB_Name                  string  `mapstructure:"db_name"`
	REDIS_Addr               string  `mapstructure:"redis_addr"`
	REDIS_Password           string  `mapstructure:"redis_password"`
	REDIS_DB                 int     `mapstructure:"redis_db"`
	AccessToken_Secret       string  `mapstructure:"access_token_secret"`
	AccessToken_ExpiryHours  int     `mapstructure:"access_token_expiry_hours"`
	RefreshToken_Secret      string  `mapstructure:"refresh_token_secret"`
	RefreshToken_ExpiryHours int     `mapstructure:"refresh_token_expiry_hours"`
	MFA_QrcodeWidth          int     `mapstructure:"mfa_qrcode_width"`
	MFA_QrcodeHeight         int     `mapstructure:"mfa_qrcode_height"`
	IO_ImageSizeLimit        float64 `mapstructure:"io_image_size_limit"` // mb
	IO_VideoSizeLimit        float64 `mapstructure:"io_video_size_limit"` // mb
	Video_DefaultPageSize    int     `mapstructure:"video_default_page_size"`
	Social_DefaultPageSize   int     `mapstructure:"social_default_page_size"`
	Etcd_Addr                string  `mapstructure:"etcd_addr"`
}

var (
	env_ *env
	once sync.Once
)

func Instance() *env {
	once.Do(func() {
		env_ = getInstance()
	})
	return env_
}

func getInstance() *env {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("err reading .env: ", err.Error())
	}

	var e env
	if err := v.Unmarshal(&e); err != nil {
		log.Fatal("err resolving .env: ", err.Error())
	}

	return &e
}
