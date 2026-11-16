package sso_grpc

import (
	"fmt"

	ssov1 "github.com/nktfmnko/sso_protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api ssov1.AuthClient
}

func New(config Config) (*Client, error) {
	maxAttempts := config.RetriesCount + 1
	retryPolicy := fmt.Sprintf(`{
		"methodConfig": [{
			"name": [{"service": ""}],
			"retryPolicy": {
				"maxAttempts": %d,
				"initialBackoff": "0.1s",
				"maxBackoff": "1s",
				"backoffMultiplier": 2.0,
				"retryableStatusCodes": ["UNAVAILABLE"]
			}
		}]
	}`, maxAttempts)

	conn, err := grpc.NewClient(config.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
	)

	if err != nil {
		return nil, fmt.Errorf("new grpc sso client: %w", err)
	}

	return &Client{api: ssov1.NewAuthClient(conn)}, nil
}
