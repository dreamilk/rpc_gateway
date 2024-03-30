package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dreamilk/rpc_gateway/config"
	"github.com/dreamilk/rpc_gateway/handler"
	"github.com/dreamilk/rpc_gateway/log"
)

func main() {
	fmt.Println("hello rpc_gataway")

	http.HandleFunc("/", handler.ServiceGateway)

	if err := http.ListenAndServe(config.DeployConf.Addr, nil); err != nil {
		log.Error(context.TODO(), "ListenAndServe failed")
	}

}
