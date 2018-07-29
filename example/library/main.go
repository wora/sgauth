package main

import "fmt"
import "context"

import (
	"os"
)

var defaultScopes = "https://www.googleapis.com/auth/xapi.zoo"

func main() {
	args, err := parseArguments()
	if err != nil {
		printUsage()
		println(err.Error())
		return
	}

	if os.Args[1] == "protorpc" {
		if len(os.Args) < 3 {
			printUsage()
			return
		}
		c, err := newHTTPClient(context.Background(), args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		protoRPCExample(c)
	} else if os.Args[1] == "grpc" {
		c, err := newGrpcClient(context.Background(), args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		gRPCExample(c)
	} else {
		printUsage()
	}
}