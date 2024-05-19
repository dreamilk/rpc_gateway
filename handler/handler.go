package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/dreamilk/rpc_gateway/log"
	"github.com/dreamilk/rpc_gateway/utils"
)

type Response struct {
	RetCode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ServiceGateway(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()
	ctx = utils.WithTraceId(ctx, uuid.NewString())

	// check url valid
	var input proto.Message
	var output proto.Message

	// invoke rpc
	url := ""
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "", zap.Error(err))
	}

	rpcPath := ""

	err = conn.Invoke(ctx, rpcPath, input, output)
	if err != nil {
		log.Error(ctx, "", zap.Error(err))
	}

	b, err := json.Marshal(output)
	if err != nil {
		log.Error(ctx, "", zap.Error(err))
	}

	// send response

	log.Infof(ctx, "host:%s %s", req.Host, utils.TraceId(ctx))
	fmt.Fprint(w, string(b))

}
