package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"go.uber.org/zap"

	"github.com/dreamilk/rpc_gateway/log"
)

type ApiInfo struct {
	AppName     string
	ServiceName string
	Method      string
	InputType   *desc.MessageDescriptor
	OuputType   *desc.MessageDescriptor
}

var m = make(map[string]*ApiInfo)
var mtx sync.Mutex

func getApiInfo(ctx context.Context, api *Api) (*ApiInfo, error) {
	mtx.Lock()
	defer mtx.Unlock()

	if _, ok := m[api.Url]; !ok {
		log.Infof(ctx, "init this url:%s", api.Url)

		filePath := "./api/" + api.AppName + ".proto"

		p := protoparse.Parser{}
		fds, err := p.ParseFiles(filePath)
		if err != nil {
			return nil, err
		}

		sd := fds[0].FindService(api.ServiceName)
		if sd == nil {
			return nil, fmt.Errorf("not found service: %s", api.ServiceName)
		}

		md := sd.FindMethodByName(api.MethodName)
		if md == nil {
			return nil, fmt.Errorf("not found method: %s", api.MethodName)
		}

		input := md.GetInputType()
		output := md.GetOutputType()

		m[api.Url] = &ApiInfo{
			AppName:     api.AppName,
			ServiceName: api.ServiceName,
			Method:      api.RpcMethod,
			InputType:   input,
			OuputType:   output,
		}

	}

	return m[api.Url], nil
}

func genProtoMessage(ctx context.Context, api *Api, b []byte) (proto.Message, proto.Message, error) {
	apiInfo, err := getApiInfo(ctx, api)
	if err != nil {
		log.Error(ctx, "no found this api", zap.Error(err))
		return nil, nil, err
	}

	dymsgInput := dynamic.NewMessage(apiInfo.InputType)
	dymsgOuput := dynamic.NewMessage(apiInfo.OuputType)

	opt := jsonpb.Unmarshaler{}
	opt.AllowUnknownFields = true

	if err := dymsgInput.UnmarshalJSONPB(&opt, b); err != nil {
		log.Error(ctx, "UnmarshalJSONPB failed", zap.Error(err))
		return nil, nil, err
	}

	return dymsgInput, dymsgOuput, nil
}
