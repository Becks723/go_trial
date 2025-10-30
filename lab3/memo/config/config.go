package config

type config struct {
	Server  serverConfig  `mapstructure:"server"`
	MySQL   mySQLConfig   `mapstructure:"mysql"`
	General generalConfig `mapstructure:"general"`
}

type serverConfig struct {
	Port int `mapstructure:"port"`
}

type mySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type generalConfig struct {
	DefaultLimit int `mapstructure:"limit"`
}
