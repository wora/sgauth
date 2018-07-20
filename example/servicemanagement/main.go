package main

import "fmt"
import "context"

import (
	"github.com/wora/protorpc/client"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/api/servicemanagement/v1"
	"github.com/shinfan/sgauth/oauth2"
	"os"
	"net/http"
	"google.golang.org/grpc"
)

func NewHTTPClient(ctx context.Context, baseUrl string, use_jwt bool, aud string) (*client.Client, error) {
	var http *http.Client
	var err error
	if (use_jwt) {
		http, err = oauth2.JWTClient(ctx, aud,"https://www.googleapis.com/auth/cloud-platform")
	} else {
		http, err = oauth2.DefaultClient(ctx, "https://www.googleapis.com/auth/cloud-platform")
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

func NewGrpcClient(ctx context.Context, use_jwt bool, aud string) (servicemanagement.ServiceManagerClient) {
	var conn *grpc.ClientConn
	if (use_jwt) {
		conn, _ = oauth2.JWTGrpcConn(ctx, "servicemanagement.googleapis.com", "443",  aud)
	} else {
		conn, _ = oauth2.DefaultGrpcConn(ctx, "servicemanagement.googleapis.com", "443", "https://www.googleapis.com/auth/cloud-platform")
	}
	return servicemanagement.NewServiceManagerClient(conn)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getAudience() string {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--jwt" {
			if (len(os.Args) < i + 2) {
				printUsage()
				return ""
			}
			return os.Args[i + 1]
		}
	}
	return ""
}


func printUsage() {
	fmt.Println("Usage: cmd [grpc|protorpc] [--jwt aud] [baseUrl]")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cmd [grpc|protorpc] [--jwt aud] [baseUrl]")
		return
	}

	use_jwt := false
	aud := ""
	if contains(os.Args, "--jwt"){
		use_jwt = true;
		aud = getAudience()
	}


	if os.Args[1] == "protorpc" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: cmd http baseUrl")
			return
		}

		baseUrl := os.Args[len(os.Args) - 1]
		c, err := NewHTTPClient(context.Background(), baseUrl, use_jwt, aud)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		request := &servicemanagement.ListServicesRequest{}
		response := &servicemanagement.ListServicesResponse{}
		err = c.Call(context.Background(), "ListServices", request, response)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(proto.MarshalTextString(response))
		}
	} else if os.Args[1] == "grpc" {
		c := NewGrpcClient(context.Background(), use_jwt, getAudience())
		request := &servicemanagement.ListServicesRequest{}
		response, err := c.ListServices(context.Background(), request)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(proto.MarshalTextString(response))
		}
	} else  {
		fmt.Println("Usage: cmd [grpc|protorpc]")
	}
}