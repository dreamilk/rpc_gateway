package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGenMsg(t *testing.T) {

	ctx := context.Background()
	api := Api{
		AppName:     "hello_world",
		ServiceName: "api.Test",
		MethodName:  "Ping",
	}

	os.Chdir("../")

	type msg struct{}

	m := msg{}
	b, _ := json.Marshal(m)

	inpput, ouput, err := genProtoMessage(ctx, &api, b)
	fmt.Println(inpput, ouput, err)
}
