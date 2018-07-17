package main

import "fmt"
import "context"

import (
	"github.com/wora/protorpc/client"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/example/library/v1"
	"github.com/shinfan/sgauth/oauth2"
	"os"
)

func NewHTTPClient(ctx context.Context, baseUrl string) (*client.Client, error) {
	http, err := oauth2.DefaultClient(ctx, "https://www.googleapis.com/auth/xapi.zoo")
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

func NewGrpcClient(ctx context.Context) (library.LibraryServiceClient) {
	conn, _ := oauth2.DefaultGrpcConn(ctx, "https://www.googleapis.com/auth/cloud-platform")
	return library.NewLibraryServiceClient(conn)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cmd [grpc|protorpc]")
		return
	}

	if os.Args[1] == "protorpc" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: cmd http baseUrl")
			return
		}
		c, err := NewHTTPClient(context.Background(), os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		request := &library.ListShelvesRequest{}
		response := &library.ListShelvesResponse{}
		err = c.Call(context.Background(), "ListShelves", request, response)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(proto.MarshalTextString(response))
		}
	} else if os.Args[1] == "grpc" {
		c := NewGrpcClient(context.Background())
		request := &library.ListShelvesRequest{}
		response, err := c.ListShelves(context.Background(), request)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(proto.MarshalTextString(response))
		}
	} else {
		fmt.Println("Usage: cmd [grpc|protorpc]")
	}
}