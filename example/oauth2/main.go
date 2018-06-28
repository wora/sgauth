package main

import "os"
import "fmt"
import "context"

import (
	"github.com/wora/protorpc/client"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/servicemanagement/v1"
	"sgauth/oauth2/google"
)

func NewClient(ctx context.Context, baseUrl string) (*client.Client, error) {
	http, err := google.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}
	c := &client.Client{
		HTTP:        http,
		BaseURL:     baseUrl,
		UserAgent:   "protorpc/0.1",
	}
	return c, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: cmd baseUrl")
		return
	}
	c, err := NewClient(context.Background(), os.Args[1])
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	request := &servicemanagement.ListServicesRequest{}
	response := &servicemanagement.ListServicesResponse{}
	err = c.Call(context.Background(), "ListServices", request, response)
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print(proto.MarshalTextString(response))
	}
}