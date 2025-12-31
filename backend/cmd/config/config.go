package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// LoadEnv は環境変数を読み取りConfig(グローバル変数)に代入する関数
// 引数にファイルを指定することによって、特定のenvファイルを読み取る
func LoadEnv(envfile ...string) error {
	// 引数が指定されていない場合は .env を読み込む
	if len(envfile) == 0 {
		envfile = []string{".env"}
	}

	// .envファイルの読み込み（存在しない場合はエラーを返さない）
	_ = godotenv.Load(envfile...)

	config := config{}

	if err := env.Parse(&config.App); err != nil {
		return err
	}

	if err := env.Parse(&config.Database); err != nil {
		return err
	}

	if err := env.Parse(&config.OAuth); err != nil {
		return err
	}

	if err := env.Parse(&config.Session); err != nil {
		return err
	}

	Config = &config

	return nil
}
