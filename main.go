package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dreamilk/rpc_gateway/config"
	"github.com/dreamilk/rpc_gateway/handler"
)

func main() {
	fmt.Println("hello rpc_gataway")

	http.HandleFunc("/", handler.ServiceGateway)

	if err := http.ListenAndServe(config.DeployConf.Addr, nil); err != nil {
		log.Fatalln(err)
	}

}
