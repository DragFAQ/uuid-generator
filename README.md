# uuid-generator

## Description
Develop an stateful application containing generated uuid hash in its memory.

Hash should be recreated every 5 minutes.

Application should contain two api servers: gRPC and http.

Each api should implement single endpoint to get actual hash string and hash generation datetime.

Cover code with unit tests where it’s needed.

This app should demonstrate coding quality, app design skills, golang best practices, etc.

It is preferable to upload the results of the work to github or some other public repo.

## Notes

Generally, this app could be contained in just one small main.go file, especially if we're certain it won't expand. However, if it becomes part of a microservice architecture, we may want to move some packages into shared libraries and reuse it in another apps.

## How to run

In project root create your own .env file (based on [env.example](env.example))

To make clean production build and run it locally:
```shell
docker-compose up --build uuid-generator
```
To make and run build for debugging:
```shell
docker-compose up --build uuid-generator-debug
```
After running, you can connect to localhost:40000 to debug it

To run tests:
```shell
go test -v -coverprofile=coverage.out ./...
```

## Use
For HTTP endpoint call http://localhost:8080

For GRPC call grpc://localhost:8090 rpc GetCurrentHash
