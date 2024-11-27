package config

import (
    "fmt"
    "sync"

    "github.com/spf13/viper"
)

type Config struct {
    APIKey string
    Model  string
}

var (
    cfg  *Config
    once sync.Once
)

func LoadConfig() (*Config, error) {
    var err error
    once.Do(func() {
        viper.SetConfigName("config")
        viper.SetConfigType("yaml")
        viper.AddConfigPath(".")
        viper.AutomaticEnv()

        if err = viper.ReadInConfig(); err != nil {
            // 設定ファイルが見つからない場合は無視
        }

        cfg = &Config{
            APIKey: viper.GetString("OPENAI_API_KEY"),
            Model:  viper.GetString("MODEL"),
        }
    })

    if cfg.APIKey == "" {
        err = fmt.Errorf("環境変数または設定ファイルに OPENAI_API_KEY が設定されていません")
    }

    if cfg.Model == "" {
        cfg.Model = "gpt-3.5-turbo" // デフォルトモデル
    }

    return cfg, err
}
