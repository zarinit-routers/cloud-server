package grpc

import (
	"log"
	"os"

	"github.com/zarinit-routers/connector-rpc/gen/connector"
	"google.golang.org/grpc"
)

const ENV_GRPC_ADDR = "CONNECTOR_GRPC_ADDR"

func getClientsRpcAddr() string {
	addr := os.Getenv(ENV_GRPC_ADDR)
	if addr == "" {
		log.Fatal("GRPC address not set", "envVariable", ENV_GRPC_ADDR)
	}
	return addr
}

func Setup() error {
	var opts []grpc.DialOption
	conn, err := grpc.NewClient(getClientsRpcAddr(), opts...)
	if err != nil {
		return err
	}
	ClientsService = connector.NewClientsServiceClient(conn)
	return nil
}

var (
	ClientsService connector.ClientsServiceClient
)
