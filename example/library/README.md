Google Authenticator Examples
-------
## Installation

To install the Google Authenticator and its examples, run:
```
go get github.com/shinfan/sgauth/example/library
```
## Credentials
Currently Google Authenticator reads the service account JSON credential file from environment path:
1) Go to the Pantheon UI ([Prod](https://pantheon.corp.google.com/)|[TestGaia](https://pantheon-testgaia.corp.google.com))
2) Enable the corresponding API if you haven't. (E.g. Service Management API in the example bew)
2) Create the service account key.
2) Download the JSON credentials.
3) Set `$GOOGLE_APPLICATION_CREDENTIALS` to the JSON path.

## Command-line Usage
The demo main has the following usage pattern:
```
go run main.go [protorpc|grpc] [--jwt aud] [baseUrl]
```
where:

`[protorpc|grpc]` is the selector between ProtobufRPC and gRPC protocols.

`[--jwt]` is the flag if you want to use client-signed JWT token without OAuth2.0. To use JWT token you need to provide an `aud` field. For more information about how to construct the `aud` field please read: [Service account authorization without OAuth](https://developers.google.com/identity/protocols/OAuth2ServiceAccount)

`[baseUrl]` is the base HTTP URL used for ProtobufRPC. You don't need to specify this field for gRPC.

## Sample Usage
Currently the Library API **only supports TestGaia**

#### ProtoRPC
```
go run main.go protorpc https://test-xxiang-library-example.sandbox.googleapis.com/\$rpc/google.example.library.v1.LibraryService/
```
#### gRPC
```
go run main.go grpc
```

#### JWT Token
To authorize with JWT token, you can specify the `--jwt` flag, for example:
```
go run main.go protorpc --jwt https://test-xxiang-library-example.sandbox.googleapis.com/google.example.library.v1.LibraryService https://test-xxiang-library-example.sandbox.googleapis.com/\$rpc/google.example.library.v1.LibraryService/
```
