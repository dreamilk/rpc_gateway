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

	cache "github.com/dreamilk/go-cache"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dreamilk/rpc_gateway/config"
	"github.com/dreamilk/rpc_gateway/log"
	"github.com/dreamilk/rpc_gateway/utils"
)

type Response struct {
	RetCode int         `json:"retcode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var serviceAddrCache = cache.New[string](5*time.Second, 3*time.Second)

func searchService(ctx context.Context, serviceName string, tag string) (string, error) {
	// use chace to find service addr in consul
	v, ok := serviceAddrCache.Get(serviceName)
	if ok {
		return v, nil
	}
	log.Warn(ctx, "no found service in cache", zap.String("serviceName", serviceName))

	// Create a Consul API client
	conf := api.DefaultConfig()
	conf.Address = config.DeployConf.Consul

	client, err := api.NewClient(conf)
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

	serviceAddrCache.Set(serviceName, addr, cache.DefaultExpiration)

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
			err = fmt.Errorf("panic: %v", r)
		}

		if err != nil {

			resp.Data = nil
			resp.Message = err.Error()
		}

		sendMsg(ctx, w, &resp)
	}()

	requestApi, err := parseUrl(req.RequestURI)
	if err != nil {
		log.Error(ctx, "parseUrl failed", zap.Error(err))
		return
	}

	body, err := getParams(req)
	if err != nil {
		return
	}

	input, output, err := genProtoMessage(ctx, requestApi, body)
	if err != nil {
		return
	}

	addr, err := searchService(ctx, requestApi.AppName, "")
	if err != nil {
		log.Error(ctx, "serach service failed", zap.Error(err))
		return
	}
	log.Info(ctx, "addr", zap.Any("addr", addr), zap.Any("url", requestApi.Url), zap.Any("rpcMethod", requestApi.RpcMethod))

	// invoke rpc
	if err = invokeRpc(ctx, addr, requestApi.RpcMethod, input, output); err != nil {
		log.Error(ctx, "invoke rpc failed", zap.Error(err))
		return
	}

	resp.Data = output
}

var connCache = cache.New[*grpc.ClientConn](time.Minute, 5*time.Second)

func invokeRpc(ctx context.Context, addr string, rpcMethod string, input any, output any) error {
	key := addr + rpcMethod

	conn, ok := connCache.Get(key)
	if !ok {
		log.Warn(ctx, "no found conn in cache", zap.String("key", key))
		connection, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error(ctx, "grpc.Dial failed", zap.Error(err))
			return err
		}

		connCache.Set(key, connection, cache.DefaultExpiration)
		conn = connection
	}
	// invoke rpc
	return conn.Invoke(ctx, rpcMethod, input, output)
}

type Api struct {
	AppName     string
	ServiceName string
	MethodName  string
	RpcMethod   string
	Url         string
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
		MethodName:  str[3],
		RpcMethod:   fmt.Sprintf("/%s.%s/%s", str[1], str[2], str[3]),
		Url:         url,
	}, nil
}
