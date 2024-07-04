package client

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "social-network-otus/internal/client/proto"
	"social-network-otus/internal/config"
	"social-network-otus/internal/logger"
)

type Client struct {
	config     *config.Config
	conn       *grpc.ClientConn
	ClientGRPC pb.DialogServiceClient
}

func New(cfg *config.Config) *Client {
	return &Client{config: cfg}
}

func InitHooks(lc fx.Lifecycle, client *Client, logger logger.LoggerInterface) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				conn, err := grpc.Dial(client.config.GRPCDialogs, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					logger.Fatal(fmt.Sprintf("failed to dial grpc %s:", client.config.GRPCDialogs), err, nil)
				}
				client.conn = conn

				client.ClientGRPC = pb.NewDialogServiceClient(conn)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := client.conn.Close()

			return err
		},
	})
}
