package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
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

func searchService(ctx context.Context, serviceName string, tag string) (string, error) {
	// Create a Consul API client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Error(ctx, "NewClient failed", zap.Error(err))
		return "", err
	}

	svc, _, err := client.Health().Service(serviceName, tag, true, nil)
	if err != nil {
		log.Error(ctx, "", zap.Error(err))
		return "", err
	}

	if len(svc) == 0 {
		return "", fmt.Errorf("not found %s in cousul", serviceName)
	}

	// TODO
	// load balance
	addr := svc[0].Node.Address + ":" + strconv.Itoa(svc[0].Service.Port)
	return addr, nil
}

func ServiceGateway(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()
	ctx = utils.WithTraceId(ctx, uuid.NewString())

	addr, err := searchService(ctx, "consul", "")
	if err != nil {
		log.Error(ctx, "", zap.Error(err))
	}
	log.Info(ctx, "addr", zap.Any("addr", addr))

	// check url valid
	var input proto.Message
	var output proto.Message

	// invoke rpc
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "grpc.Dial failed", zap.Error(err))
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
