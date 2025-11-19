package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/laoggu/cabbage-vehicle-backend/api"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/crypt"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/kafka"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/mysql"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/oss"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/redis"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/config"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/jwt"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/logger"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/otel"
	"github.com/laoggu/cabbage-vehicle-backend/internal/service/auth"
	"github.com/laoggu/cabbage-vehicle-backend/internal/service/vehicle"
	vrepo "github.com/laoggu/cabbage-vehicle-backend/internal/service/vehicle/repo"
)

func main() {
	logger.Init()
	defer logger.Sync()
	config.Load()
	crypt.Init(config.C.SignSecret)
	jwt.Init(config.C.SignSecret)

	shut := otel.InitTracer("gateway", logger.L)
	defer shut()

	// infra
	db, err := mysql.New(config.C.MysqlDSN, logger.L)
	if err != nil {
		logger.L.Fatal("mysql", zap.Error(err))
	}
	rdb := redis.New(config.C.RedisAddr)
	kfk := kafka.NewWriter(config.C.KafkaBrokers, "RouteEvent")
	ossClient, err := oss.New(config.C.OssEndpoint, config.C.OssAK, config.C.OssSK, "cabbage-vehicle")
	if err != nil {
		logger.L.Fatal("oss", zap.Error(err))
	}

	// 内部 gRPC 服务
	grpcS := grpc.NewServer()
	api.RegisterAuthServiceServer(grpcS, auth.NewServer(logger.L))
	api.RegisterVehicleServiceServer(grpcS, vehicle.NewServer(logger.L, vrepo.NewVehicleRepo(db), ossClient))

	// gRPC-Gateway mux
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := api.RegisterAuthServiceHandlerFromEndpoint(context.Background(), gwMux, "localhost:9090", opts); err != nil {
		logger.L.Fatal("gw auth", zap.Error(err))
	}
	if err := api.RegisterVehicleServiceHandlerFromEndpoint(context.Background(), gwMux, "localhost:9090", opts); err != nil {
		logger.L.Fatal("gw vehicle", zap.Error(err))
	}

	// 健康检查
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.Handle("/", gwMux)

	// 同一端口支持 HTTP/1 & HTTP/2
	h2s := &http2.Server{}
	srv := &http.Server{
		Addr:    ":" + config.C.HttpPort,
		Handler: h2c.NewHandler(withIdempotency(mux), h2s),
	}

	// 优雅退出
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("gateway listen", zap.Error(err))
		}
	}()
	logger.L.Info("gateway listen", zap.String("port", config.C.HttpPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
