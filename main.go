package main

import (
	"github.com/DragFAQ/uuid-generator/config"
	generator "github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

var currentHash generator.Hash
var hashLock sync.RWMutex

func failOnError(logger log.Logger, err error, msg string) {
	if err != nil {
		logger.Panicf("%s: %s", msg, err)
	}
}

func main() {
	var conf = config.NewConfig()

	var logger, err = log.NewZapLogger(conf.Logger)
	failOnError(logger, err, "new logger")

	logger.Infof("starting service on %v", time.Now())

	// Start generating the initial hash
	currentHash = generator.Hash{
		Value:          uuid.New().String(),
		GenerationTime: time.Now(),
	}
	go generator.GenerateHash(&currentHash, &hashLock, logger, conf.Settings.HashTTLSeconds)

	// Define the HTTP API routes
	httpHandler := handler.NewHttpHandler(&currentHash, &hashLock, logger)
	http.HandleFunc("/", httpHandler.GetCurrentHash)

	// Define the gRPC server
	grpcServer := grpc.NewServer()
	grpcHandler := handler.NewGrpcHandler(&currentHash, &hashLock, logger)
	pb.RegisterHashServiceServer(grpcServer, grpcHandler)

	// Start both servers
	go func() {
		err = http.ListenAndServe(":"+conf.Server.HttpPort, nil)
		failOnError(logger, err, "Failed to start HTTP server")
	}()

	go func() {
		listener, err := net.Listen("tcp", ":"+conf.Server.GrpcPort)
		failOnError(logger, err, "Failed to listen port")

		err = grpcServer.Serve(listener)
		failOnError(logger, err, "Failed to start GRPC server")
	}()

	select {}
}
