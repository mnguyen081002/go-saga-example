package config

import (
	"fmt"
	"item-service/utils/constants"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	ConfigDefaultFile = "config/config.yml"
	ConfigReleaseFile = "config/config.release.yml"
	ConfigDevFile     = "config/config.dev.yml"
	configType        = "yml"
)

type (
	Config struct {
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Services       Services `mapstructure:"services"`
		Database       Database `mapstructure:"database"`
		Logger         Logger   `mapstructure:"logger"`
	}

	Server struct {
		Host     string `mapstructure:"host"`
		Env      string `mapstructure:"env"`
		UseRedis bool   `mapstructure:"useRedis"`
		Port     int    `mapstructure:"port"`
	}

	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	}

	Logger struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
		Prefix string `mapstructure:"prefix"`
	}

	Gin struct {
		Mode string `mapstructure:"mode"`
	}

	Services struct {
		PaymentServiceURL string `mapstructure:"paymentServiceURL"`
		ItemServiceURL    string `mapstructure:"itemServiceURL"`
	}
)

func NewConfig() Config {
	InitConfig()
	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable decode into config struct, %v", err)
	}
	return *conf
}

func InitConfig() {
	var configFile string
	switch os.Getenv("ENV") {
	case constants.Local:
		configFile = ConfigDefaultFile
		fmt.Printf("gin mode %s,config%s\n", gin.EnvGinMode, ConfigDefaultFile)
	case constants.Prod:
		configFile = ConfigReleaseFile
		fmt.Printf("gin mode %s,config%s\n", gin.EnvGinMode, ConfigReleaseFile)
	case constants.Dev:
		configFile = ConfigDevFile
		fmt.Printf("gin mode: %s,config file: %s\n", gin.EnvGinMode, ConfigDevFile)
	}
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}
