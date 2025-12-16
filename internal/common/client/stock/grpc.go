package stock

import (
	"context"
	"errors"
	"time"

	clipkg "github.com/loveRyujin/gorder/common/client/pkg"
	"github.com/loveRyujin/gorder/common/discovery"
	"github.com/loveRyujin/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCli(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !waitForStock(viper.GetDuration("grpc-dial-timeout") * time.Second) {
		return nil, func() error { return nil }, errors.New("grpc service is not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(_ string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewServerHandler()),
	}, nil
}

func waitForStock(timeout time.Duration) bool {
	logrus.Infof("waiting for stock grpc client, timeout: %v seconds", timeout.Seconds())
	return clipkg.WaitFor(viper.GetString("stock.grpc-addr"), timeout)
}
