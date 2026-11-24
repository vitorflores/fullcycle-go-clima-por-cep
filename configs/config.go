package configs

import (
	"errors"
	"io/fs"

	"github.com/spf13/viper"
)

type Conf struct {
	Porta string `mapstructure:"PORT"`
}

func CarregarConfig() (*Conf, error) {
	var cfg *Conf
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.SetDefault("PORT", "8080")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
