package main

import (
	"fmt"
	"log"
	"sync"

	staff2 "github.com/wayla99/go_gorm.git/src/entity/staff"
	"github.com/wayla99/go_gorm.git/src/repository/staff"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wayla99/go_gorm.git/src/interface/fiber_server"
	"github.com/wayla99/go_gorm.git/src/use_case"

	"go.uber.org/zap/zapcore"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger

type config struct {
	AppName            string `env:"APP_NAME" envDefault:"staff-test"`
	AppVersion         string `env:"APP_VERSION" envDefault:"v0.0.0"`
	Environment        string `env:"ENVIRONMENT" envDefault:"development"`
	Port               uint   `env:"PORT" envDefault:"9000"`
	Debuglog           bool   `env:"DEBUG_LOG" envDefault:"false"`
	PostgresDBEndpoint string `env:"POSTGRES_DB_ENDPOINT" envDefault:"postgres://postgres:postgres@localhost:5432"`
	PgSSL              string `env:"PG_SSL" envDefault:"disable"`
	PgDBName           string `env:"PG_DB_NAME" envDefault:"staff"`
	PgDBTable          string `env:"PG_DB_TABLE" envDefault:"staff"`
}

func main() {
	cfg := initEnvironment()
	initLogger(cfg)

	staffRepo := initRepositories(cfg)
	useCase := use_case.New(staffRepo)
	initInterfaces(cfg, useCase)

	//user, err := query.Staff.Where(query.Staff.FirstName.Eq("string")).First()
	//if err != nil {
	//	log.Println("err")
	//}
	//log.Println("uiserdsfsdfsef : ", user)
}

func initEnvironment() config {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}

	var cfg config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Error parse env: %s\n", err)
	}

	return cfg
}

func initLogger(cfg config) {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	LogLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	if cfg.Debuglog {
		LogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	conf.Level = LogLevel

	lg, err := conf.Build()
	if err != nil {
		log.Fatalf("Error build logger: %s\n", err)
	}
	defer lg.Sync()

	zap.ReplaceGlobals(lg)
	logger = zap.L().Named("bootstrap")
	logger.Info("Logger initialized")
}

func initRepositories(cfg config) use_case.StaffRepository {
	staffRepo, err := staff.New(cfg.PostgresDBEndpoint, cfg.PgDBName, cfg.PgDBTable, cfg.PgSSL, &staff2.Staff{})
	if err != nil {
		logger.Fatal("Error init staff repository", zap.Error(err))
	}
	logger.Info("Staff repository initialized")
	logger.Info("Repositories initialized")

	return staffRepo
}

func initInterfaces(cfg config, useCase *use_case.UseCase) {
	wg := new(sync.WaitGroup)
	prom := prometheus.NewRegistry()

	serv := fiber_server.New(useCase, prom, &fiber_server.ServerConfig{
		AppVersion:    cfg.AppVersion,
		ListenAddress: fmt.Sprintf(":%d", cfg.Port),
		RequestLog:    true,
	})
	logger.Info("Fiber server initialized")

	serv.Start(wg)
	logger.Info("Fiber server started")

	wg.Wait()
	logger.Info("Application stopped")
}
