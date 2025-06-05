package config

import (
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/spf13/viper"
)

var (
	configFiles = []string{"database", "redis", "server"}
	Conf        = &Config{}
)

func init() {
	v := viper.New()

	for _, cf := range configFiles {
		v.SetConfigName(cf)
		v.SetConfigType("yaml")
		v.AddConfigPath("./config")

		if err := v.MergeInConfig(); err != nil {
			// 若文件不存在则跳过错误
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Logger.Fatalf("Error reading config file %s.yaml, %v", cf, err)
			}
		}
	}
	if err := v.Unmarshal(Conf); err != nil {
		log.Logger.Fatalf("Unable to decode into struct, %v", err)
	}
	log.Logger.Println("初始化配置文件...success")
}
