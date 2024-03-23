package handler

import (
	"fmt"
	"net/http"
)

func ServiceGateway(w http.ResponseWriter, req *http.Request) {

	fmt.Println(req.RequestURI)

	fmt.Fprint(w, "hello gw")
}
