package main

import (
	"github.com/wora/protorpc/client"
	"fmt"
	"google.golang.org/genproto/googleapis/example/library/v1"
	"github.com/shinfan/sgauth"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

func createSettings(args map[string]string) (*sgauth.Settings) {
	if args[kApiKey] != "" {
		return &sgauth.Settings{
			APIKey: args[kApiKey],
		}
	} else if args[kAud] != "" {
		return &sgauth.Settings{
			Audience: args[kAud],
		}
	} else {
		return &sgauth.Settings{
			Scope: args[kScope],
		}
	}
}

func newHTTPClient(ctx context.Context, args map[string]string) (
	*client.Client, error) {
	baseUrl := fmt.Sprintf("https://%s/$rpc/%s/", args[kHost], args[kApiName])

	http, err := sgauth.NewHTTPClient(ctx, createSettings(args))
	if err != nil {
		return nil, err
	}
	return &client.Client{
		HTTP:        http,
		BaseURL:     baseUrl,
		UserAgent:   "protorpc/0.1",
	}, nil
}

func newGrpcClient(ctx context.Context, args map[string]string) (library.LibraryServiceClient, error) {
	conn, err := sgauth.NewGrpcConn(ctx, createSettings(args), args[kHost], "443")
	if err != nil {
		return nil, err
	}
	return library.NewLibraryServiceClient(conn), nil
}

func protoRPCExample(client *client.Client) {
	request := &library.ListShelvesRequest{}
	response := &library.ListShelvesResponse{}
	err := client.Call(context.Background(), "ListShelves", request, response)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(proto.MarshalTextString(response))
	}
}

func gRPCExample(client library.LibraryServiceClient) {
	request := &library.ListShelvesRequest{}
	response, err := client.ListShelves(context.Background(), request)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(proto.MarshalTextString(response))
	}
}
