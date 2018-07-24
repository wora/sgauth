package main

import "fmt"
import "context"

import (
	"github.com/wora/protorpc/client"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/example/library/v1"
	"github.com/shinfan/sgauth"
	"os"
	"google.golang.org/grpc"
)

func NewHTTPClient(ctx context.Context, service_name string, api_name string, use_jwt bool) (*client.Client, error) {
	var credentials = &sgauth.Credentials{
		ServiceAccount: &sgauth.ServiceAccount{
			EnableOAuth: !use_jwt,
			ServiceName: service_name,
			APIName: api_name,
			Scopes: []string{"https://www.googleapis.com/auth/xapi.zoo"},
		},
	}
	return client.NewClient(ctx, credentials)
}

func NewGrpcClient(ctx context.Context, service_name string, use_jwt bool, aud string) (library.LibraryServiceClient) {
	var conn *grpc.ClientConn
	if (use_jwt) {
		conn, _ = sgauth.JWTGrpcConn(ctx, service_name, "443",  aud)
	} else {
		conn, _ = sgauth.DefaultGrpcConn(ctx,
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
	} else {
		println("Error: --service_name is required")
		printUsage()
		return
	}

	if contains(os.Args, "--jwt") {
		use_jwt = true;
		aud = fmt.Sprintf("https://%s/%s", service_name, getFlagValue("--api_name"))
		println(aud)
	}

	if contains(os.Args, "--api_name") {
		api_name = getFlagValue("--api_name")
	} else if (use_jwt || os.Args[1] == "protorpc") {
		println("Error: --api_name is required in JWT or ProtoRPC mode")
		printUsage()
		return
	}

	if os.Args[1] == "protorpc" {
		if len(os.Args) < 3 {
			printUsage()
			return
		}

		c, err := NewHTTPClient(context.Background(), service_name, api_name, use_jwt)
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