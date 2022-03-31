package main

import (
	"github.com/byteintellect/go_commons"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/go_commons/entity"
	"github.com/byteintellect/otp_svc/config"
	"github.com/byteintellect/otp_svc/pkg/domain"
	"github.com/byteintellect/otp_svc/pkg/svc"
	otpsv1 "github.com/byteintellect/protos_go/otps/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"time"
)

func NewGrpcServer(cfg *config.OtpServiceConfig, app *go_commons.BaseApp) (*grpc.Server, error) {
	grpcMux := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Duration(cfg.BaseConfig.ServerConfig.KeepAliveTime) * time.Second,
				Timeout: time.Duration(cfg.BaseConfig.ServerConfig.KeepAliveTimeOut) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// logging middleware
				grpcZap.UnaryServerInterceptor(app.Logger()),

				// Request-Id interceptor
				requestid.UnaryServerInterceptor(),

				// Metrics middleware
				app.GrpcMetrics().UnaryServerInterceptor(),

				// validation middleware
				grpcValidator.UnaryServerInterceptor(),

				// collection operators middleware
				gateway.UnaryServerInterceptor(),

				// trace middleware
				otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(app.Tracer())),
			),
		),
		grpc.StreamInterceptor(app.GrpcMetrics().StreamServerInterceptor()),
	)
	domainMappings := getDomainMappings()
	otpSvc := svc.NewOtpSvc(db.NewGORMRepository(db.WithDb(app.Db()), db.WithCreator(domainMappings.GetMapping("otps"))), cfg, app.Logger())
	otpsv1.RegisterOtpServiceServer(grpcMux, otpSvc)
	// Register reflection service on gRPC server.
	reflection.Register(grpcMux)
	grpcPrometheus.Register(grpcMux)
	app.GrpcMetrics().InitializeMetrics(grpcMux)
	return grpcMux, nil
}

func getDomainMappings() entity.DomainFactory {
	factory := entity.NewDomainFactory()
	factory.RegisterMapping("otps", func() entity.Base {
		return &domain.Otp{}
	})
	return *factory
}
