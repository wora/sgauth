package main

import "fmt"
import "context"

import (
	"github.com/wora/protorpc/client"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/example/library/v1"
	"github.com/shinfan/sgauth/oauth2"
	"os"
	"net/http"
	"google.golang.org/grpc"
)

func NewHTTPClient(ctx context.Context, baseUrl string, use_jwt bool, aud string) (*client.Client, error) {
	var http *http.Client
	var err error
	if (use_jwt) {
		http, err = oauth2.JWTClient(ctx, aud,"https://www.googleapis.com/auth/xapi.zoo")
	} else {
		http, err = oauth2.DefaultClient(ctx, "https://www.googleapis.com/auth/xapi.zoo")
	}
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

func NewGrpcClient(ctx context.Context, service_name string, use_jwt bool, aud string) (library.LibraryServiceClient) {
	var conn *grpc.ClientConn
	if (use_jwt) {
		conn, _ = oauth2.JWTGrpcConn(ctx, service_name, "443",  aud)
	} else {
		conn, _ = oauth2.DefaultGrpcConn(ctx,
			service_name, "443", "https://www.googleapis.com/auth/xapi.zoo")
	}
	return library.NewLibraryServiceClient(conn)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getFlagValue(flag string) string {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == flag {
			if (len(os.Args) < i + 2) {
				printUsage()
				return ""
			}
			return os.Args[i + 1]
		}
	}
	printUsage()
	return ""
}


func printUsage() {
	fmt.Println("Usage: cmd [grpc|protorpc] [--jwt] [--service_name] [--api_name]")
}

func main() {
	if len(os.Args) < 3 {
		printUsage()
		return
	}

	use_jwt := false
	aud := ""
	service_name := ""
	api_name := ""

	if contains(os.Args, "--service_name") {
		service_name = getFlagValue("--service_name")
	}

	if contains(os.Args, "--api_name") {
		api_name = getFlagValue("--api_name")
	}

	if contains(os.Args, "--jwt") {
		use_jwt = true;
		aud = fmt.Sprintf("https://%s/%s", service_name, api_name)
	}

	if os.Args[1] == "protorpc" {
		if len(os.Args) < 3 {
			printUsage()
			return
		}

		baseUrl := fmt.Sprintf("https://%s/$rpc/%s/", service_name, api_name)
		c, err := NewHTTPClient(context.Background(), baseUrl, use_jwt, aud)
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
		c := NewGrpcClient(context.Background(), service_name, use_jwt, aud)
		request := &library.ListShelvesRequest{}
		response, err := c.ListShelves(context.Background(), request)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(proto.MarshalTextString(response))
		}
	} else {
		printUsage()
	}
}