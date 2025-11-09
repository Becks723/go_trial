package config

type config struct {
	Server  serverConfig  `mapstructure:"server"`
	MySql   mySqlConfig   `mapstructure:"mysql"`
	General generalConfig `mapstructure:"general"`
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
