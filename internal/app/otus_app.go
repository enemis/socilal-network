package app

import (
	"go.uber.org/fx"
	"social-network-otus/internal/cache"
	"social-network-otus/internal/queue"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/config"
	"social-network-otus/internal/database"
	"social-network-otus/internal/friend"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/post"
	"social-network-otus/internal/rest"
	"social-network-otus/internal/session"
	"social-network-otus/internal/token"
	"social-network-otus/internal/user"
)

type App struct {
	container *fx.App
}

func NewApp() (*App, error) {
	//ctx := context.Background()

	fxContainer := fx.New(
		config.Module,
		cache.Module,
		logger.Module,
		database.Module,
		rest.Module,
		auth.Module,
		session.Module,
		token.Module,
		user.Module,
		friend.Module,
		post.Module,
		queue.Module,
	)
	return &App{
		container: fxContainer,
	}, nil
}

func (app *App) Run() error {
	app.container.Run()
	return nil
}

// 	logrus.SetFormatter(new(logrus.JSONFormatter))
// 	if err := initConfig(); err != nil {
// 		logrus.Fatalf("error init config %s", err.Error())
// 	}

// 	if err := gotenv.Load(); err != nil {
// 		logrus.Fatalf("error load .env: %s", err.Error())
// 	}

// 	db, err := repository.NewPostgresDb(repository.Config{
// 		Host:     viper.GetString("db_host"),
// 		Port:     viper.GetString("db_port"),
// 		Username: viper.GetString("db_username"),
// 		DBName:   viper.GetString("db_name"),
// 		SSLMode:  viper.GetString("db_ssl_mode"),
// 		Password: os.Getenv("DB_PASS"),
// 	})

// 	if err != nil {
// 		log.Fatalf("failed init db:%s", err.Error())
// 	}

// 	repos := repository.NewRepository(db)
// 	services := service.NewService(repos)
// 	handlers := handler.NewHandler(services)
// 	srv := new(server.Server)

// 	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
// 		log.Fatalf("error occures while running http server: %s", err.Error())
// 	}
// }

// func initConfig() error {
// 	viper.AddConfigPath("configs")
// 	viper.SetConfigName("config")
// 	return viper.ReadInConfig()
// }
