package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	svc, _, err := client.Health().Service(serviceName, tag, false, nil)
	if err != nil {
		log.Error(ctx, "search service failed", zap.Error(err))
		return "", err
	}

	if len(svc) == 0 {
		return "", fmt.Errorf("rpc service %s unavailable", serviceName)
	}

	// TODO
	// load balance
	addr := svc[0].Node.Address + ":" + strconv.Itoa(svc[0].Service.Port)
	return addr, nil
}

func ServiceGateway(w http.ResponseWriter, req *http.Request) {
	var err error
	var resp Response

	ctx, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()
	ctx = utils.WithTraceId(ctx, uuid.NewString())

	defer func() {
		if r := recover(); r != nil {
			log.Errorf(ctx, "panic %v", r)
		}

		if err != nil {

			resp.Data = nil
			resp.Message = err.Error()
		}

		sendMsg(ctx, w, &resp)
	}()

	api, err := parseUrl(req.RequestURI)
	if err != nil {
		log.Error(ctx, "parseUrl failed", zap.Error(err))
		return
	}
	log.Info(ctx, "parse result", zap.Any("api", api))

	body, err := getParams(req)
	if err != nil {
		return
	}

	input, output, err := genProtoMessage(ctx, api, body)
	if err != nil {
		return
	}

	addr, err := searchService(ctx, api.AppName, "")
	if err != nil {
		log.Error(ctx, "serach service failed", zap.Error(err))
		return
	}
	log.Info(ctx, "addr", zap.Any("addr", addr))

	// invoke rpc
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "grpc.Dial failed", zap.Error(err))
		return
	}

	rpcPath := api.RpcMethod

	err = conn.Invoke(ctx, rpcPath, input, output)
	if err != nil {
		log.Error(ctx, "conn invoke failed", zap.Error(err))
		return
	}

	resp.Data = output
}

type Api struct {
	AppName     string
	ServiceName string
	Path        string
	RpcMethod   string
}

func getParams(req *http.Request) ([]byte, error) {
	if req.Method == "GET" {
		q := req.URL.Query()
		return json.Marshal(q)
	} else if req.Method == "POST" {
		return io.ReadAll(req.Body)
	}
	return nil, fmt.Errorf("unsupport method%s", req.Method)
}

func genProtoMessage(ctx context.Context, api *Api, b []byte) (proto.Message, proto.Message, error) {
	filePath := "./api/" + api.AppName + ".proto"

	p := protoparse.Parser{}
	fds, err := p.ParseFiles(filePath)
	if err != nil {
		return nil, nil, err
	}

	sd := fds[0].FindService(api.ServiceName)
	if sd == nil {
		return nil, nil, fmt.Errorf("not found service: %s", api.ServiceName)
	}

	md := sd.FindMethodByName(api.Path)
	if md == nil {
		return nil, nil, fmt.Errorf("not found method: %s", api.Path)
	}

	input := md.GetInputType()
	output := md.GetOutputType()

	dymsgInput := dynamic.NewMessage(input)
	dymsgOuput := dynamic.NewMessage(output)

	opt := jsonpb.Unmarshaler{}
	opt.AllowUnknownFields = true

	if err := dymsgInput.UnmarshalJSONPB(&opt, b); err != nil {
		log.Error(ctx, "UnmarshalJSONPB failed", zap.Error(err))
		return nil, nil, err
	}

	return dymsgInput, dymsgOuput, nil
}

func sendMsg(ctx context.Context, w http.ResponseWriter, resp *Response) error {
	b, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "json marshal failed", zap.Error(err))
		return err
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(b))

	return nil
}

func parseUrl(url string) (*Api, error) {
	url = strings.Split(url, "?")[0]
	str := strings.Split(url, "/")
	if len(str) != 4 {
		return nil, fmt.Errorf("path:%s parse failed", url)
	}

	return &Api{
		AppName:     str[1],
		ServiceName: str[2],
		Path:        str[3],
		RpcMethod:   str[2] + "/" + str[3],
	}, nil
}
