package main

import (
	"fmt"
	"log"
	"net/http"

	"gateway/config"
	"gateway/handler"
)

func main() {
	fmt.Println("hello rpc_gataway")

	http.HandleFunc("/", handler.ServiceGateway)

	if err := http.ListenAndServe(config.DeployConf.Addr, nil); err != nil {
		log.Fatalln(err)
	}

}
