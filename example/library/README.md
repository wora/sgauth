Google Authenticator(GA) Prototype
-------
## Installation

To install the Google Authenticator prototype, run:
```
go get github.com/shinfan/sgauth/example/library
```
## Credentials
Currently Google Authenticator reads the service account JSON credential file from environment path:
1) Go to the Pantheon UI ([Prod](https://pantheon.corp.google.com/)|[TestGaia](https://pantheon-testgaia.corp.google.com))
2) Enable the corresponding API if you haven't. (E.g. Service Management API in the example below)
2) Create the service account key.
2) Download the JSON credentials.
3) Set `$GOOGLE_APPLICATION_CREDENTIALS` to the JSON path.

## Command-line Usage
The demo main has the following usage pattern:
```
go run main.go protorpc|grpc [--jwt] --service_name {service_name} [--api_name {api_name}]
```
where:

- `protorpc|grpc` *[REQUIRED]* is the selector between ProtobufRPC and gRPC protocols. 
- `[--jwt]` is the flag if you want to use client-signed JWT token without OAuth2.0. For more information about JWT token please read: [Service account authorization without OAuth](https://developers.google.com/identity/protocols/OAuth2ServiceAccount)
- `[--service_name]` *[REQUIRED]* is the full host name of the API service. e.g. test-xxiang-library-example.sandbox.googleapis.com 
- `[--api_name]` is the full API name. e.g. google.example.library.v1.LibraryService. Tjos field is only required when `protorpc` or `jwt` option is selected.

## Sample Usage

### Work with Test GAIA

The following commands run the example with the Test GAIA instance so that your credential JSON needs to be generated from Test GAIA Pantheon. Currently the API service is hosted within a sandbox environment for prototyping purpose.

#### ProtoRPC
```
go run main.go protorpc --service_name test-xxiang-library-example.sandbox.googleapis.com --api_name google.example.library.v1.LibraryService
```
#### gRPC
```
go run main.go grpc --service_name test-xxiang-library-example.sandbox.googleapis.com
```
#### JWT Token
To authorize with JWT token, you only need specify the extra `--jwt` flag, for example:
```
go run main.go grpc --jwt --service_name test-xxiang-library-example.sandbox.googleapis.com -api_name google.example.library.v1.LibraryService
```
Note: `api_name` is required for `jwt` option.

### Work with Prod GAIA

If you want to work with Prod GAIA, you can switch to use the public Library API service and everything else should be the same. e.g.
```
go run main.go grpc --service_name library-example.googleapis.com
```
