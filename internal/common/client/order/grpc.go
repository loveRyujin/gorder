package order

import (
	"context"
	"errors"
	"time"

	clipkg "github.com/loveRyujin/gorder/common/client/pkg"
	"github.com/loveRyujin/gorder/common/discovery"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCli(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !waitForOrder(viper.GetDuration("grpc-dial-timeout") * time.Second) {
		return nil, func() error { return nil }, errors.New("order service is not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for order grpc")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(_ string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}

func waitForOrder(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc client, timeout: %v seconds", timeout.Seconds())
	return clipkg.WaitFor(viper.GetString("order.grpc-addr"), timeout)
}
