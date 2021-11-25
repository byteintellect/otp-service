package main

import (
	"github.com/byteintellect/go_commons"
	"github.com/byteintellect/go_commons/config"
	config2 "github.com/byteintellect/otp_svc/config"
	otpsv1 "github.com/byteintellect/protos_go/otps/v1"
	"github.com/infobloxopen/atlas-app-toolkit/gorm/resource"
	"go.uber.org/zap"
	"log"
)

func getConfig() *config2.OtpServiceConfig {
	var cfg config2.OtpServiceConfig
	err := config.ReadFile("CONFIG_PATH", &cfg)
	if err != nil {
		log.Fatalf("error reading config")
	}
	return &cfg
}

func main() {
	cfg := getConfig()
	app, err := go_commons.NewBaseApp(&cfg.BaseConfig)
	if err != nil {
		log.Fatalf("Error initializing application %v", err)
	}
	grpcServer, err := NewGrpcServer(cfg, app)
	if err != nil {
		log.Fatalf("Error initializing gRPC server %v", err)
	}
	doneC := make(chan error)

	// Init External
	go func() {
		doneC <- go_commons.ServeExternal(&cfg.BaseConfig, app, grpcServer, otpsv1.RegisterOtpServiceHandlerFromEndpoint)
	}()
	if err := <-doneC; err != nil {
		app.Logger().Fatal("Error Starting gRPC service", zap.Error(err))
	}
	resource.RegisterApplication(cfg.BaseConfig.AppName)
	resource.SetPlural()
}
