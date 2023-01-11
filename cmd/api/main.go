package main

import (
	"os"
	"sqlmaster/internal/repository"
	"sqlmaster/internal/repository/dbrepo"
	"sqlmaster/internal/structs"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	Host        string `mapstructure:"HOST"`
	Port        int    `mapstructure:"PORT"`
	DatabaseDSN string `mapstructure:"DSN"`
	CsrfKey     string `mapstructure:"CSRF_KEY"`
	LogPath     string `mapstructure:"LOG_PATH"`
}

func getLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// The format time can be customized
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func main() {
	var (
		cfg config
		err error
	)

	MANDATORY_ENV_VARS := map[string]string{
		"HOST":     "HOST environment variable needs to be set",
		"PORT":     "PORT environment variable needs to be set",
		"DSN":      "DSN environment variable needs to be set",
		"CSRF_KEY": "CSRF_KEY environment variable needs to be set",
	}

	viper.AutomaticEnv()

	for k := range MANDATORY_ENV_VARS {
		if !viper.IsSet(k) {
			panic(MANDATORY_ENV_VARS[k])
		}
	}

	//Bind env vars
	for _, k := range []string{
		"HOST", "PORT", "DSN", "CSRF_KEY", "LOG_PATH"} {
		viper.BindEnv(k, k)
	}

	// hooks := mapstructure.ComposeDecodeHookFunc(
	// 	mapstructure.StringToTimeDurationHookFunc(),
	// 	mapstructure.StringToSliceHookFunc(","),
	// )

	cfg.LogPath = "./logs"

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	// log.Printf("Config ===> %+v\n", cfg)

	logfile := structs.LogFile{
		LogFileName: "sqlmaster",
		LogPath:     cfg.LogPath,
	}

	syncer := zap.CombineWriteSyncers(os.Stdout, &logfile)
	encoder := getLogEncoder()
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	// Print function lines
	applogger := zap.New(core, zap.AddCaller()).Sugar()

	app := application{
		Host:        cfg.Host,
		Port:        cfg.Port,
		DatabaseDSN: cfg.DatabaseDSN,
		Logger:      applogger,
	}

	conn, err := app.ConnectToDB()
	if err != nil {
		app.Logger.Fatal("Cannot connect to DB: ", err)
	}

	app.DB = &dbrepo.MysqlDBRepo{DB: conn, Timeout: 30 * time.Second}
	app.Models.Datasources = repository.NewDatasourceManager(app.DB)
	app.Models.SQLQueries = repository.NewSQLQueryManager(app.DB)
	app.Models.SQLRequests = repository.NewSQLRequestsManager(app.DB)

	defer app.DB.Connection().Close()
	app.Run()

}
