package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/dreamilk/rpc_gateway/log"
	"github.com/dreamilk/rpc_gateway/utils"
)

func ServiceGateway(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = utils.WithTraceId(ctx, uuid.NewString())

	log.Infof(ctx, "host:%s %s", req.Host, utils.TraceId(ctx))

	fmt.Fprint(w, "hello gw")
}
