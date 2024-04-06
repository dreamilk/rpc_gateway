package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"go.uber.org/zap"

	"github.com/dreamilk/rpc_gateway/config"
	"github.com/dreamilk/rpc_gateway/handler"
	"github.com/dreamilk/rpc_gateway/log"
)

func main() {
	ctx := context.Background()
	// Create a Consul API client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Error(ctx, "NewClient failed", zap.Error(err))
		return
	}
	// Create an instance representing this service. "my-service" is the
	// name of _this_ service. The service should be cleaned up via Close.
	svc, err := connect.NewService("gateway", client)
	if err != nil {
		log.Error(ctx, "new service failed", zap.Error(err))
		return
	}
	defer svc.Close()

	fmt.Println("hello rpc_gataway")

	http.HandleFunc("/", handler.ServiceGateway)

	if err := http.ListenAndServe(config.DeployConf.Addr, nil); err != nil {
		log.Error(context.TODO(), "ListenAndServe failed", zap.Error(err))
	}

}
