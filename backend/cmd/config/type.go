package config

var Config *config

type config struct {
	App struct {
		Port        string `env:"PORT" envDefault:"8080"`
		FrontendURL string `env:"FRONTEND_URL" envDefault:"http://localhost:5173"`
	}

	Database struct {
		// DATABASE_URL（Railway等で使用する一括設定）
		// 例: postgresql://user:password@host:port/dbname?sslmode=require
		URL      string `env:"DATABASE_URL"`
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		Name     string `env:"DB_NAME" envDefault:"todoapp"`
		SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	}

	OAuth struct {
		Google struct {
			ClientID     string `env:"GOOGLE_CLIENT_ID"`
			ClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
			RedirectURL  string `env:"GOOGLE_REDIRECT_URL" envDefault:"http://localhost:8080/auth/google/callback"`
		}
		Github struct {
			ClientID     string `env:"GITHUB_CLIENT_ID"`
			ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
			RedirectURL  string `env:"GITHUB_REDIRECT_URL" envDefault:"http://localhost:8080/auth/github/callback"`
		}
	}

	Session struct {
		Secret string `env:"SESSION_SECRET" envDefault:"your-secret-key-change-in-production"`
	}
}
