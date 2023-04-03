package main

import (
	"github.com/DragFAQ/uuid-generator/config"
	generator "github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var currentHash generator.Hash
var hashLock sync.RWMutex

var (
	shutDownCh     = make(chan os.Signal, 1)
	mainShutCh     = make(chan os.Signal, 1)
	generateShutCh = make(chan os.Signal, 1)
)

func failOnError(logger log.Logger, err error, msg string) {
	if err != nil && err != http.ErrServerClosed {
		logger.Panicf("%s: %s", msg, err)
	}
}

func SetupShutdownHardware() sync.WaitGroup {
	signal.Notify(shutDownCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		for sig := range shutDownCh {
			mainShutCh <- sig
			generateShutCh <- sig
		}
	}()
	return sync.WaitGroup{}
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

	wg := SetupShutdownHardware()

	// Launch goroutine with worker function
	wg.Add(1)
	go generator.GenerateHash(&currentHash, &hashLock, logger, conf.Settings.HashTTLSeconds, generateShutCh, &wg)

	// Define the HTTP API routes
	httpHandler := handler.NewHttpHandler(&currentHash, &hashLock, logger)
	http.HandleFunc("/", httpHandler.GetCurrentHash)
	srv := &http.Server{Addr: ":" + conf.Server.HttpPort}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		failOnError(logger, err, "Failed to start HTTP server")
	}()

	// Define the gRPC server
	grpcServer := grpc.NewServer()
	grpcHandler := handler.NewGrpcHandler(&currentHash, &hashLock, logger)
	pb.RegisterHashServiceServer(grpcServer, grpcHandler)

	wg.Add(1)
	go func() {
		defer wg.Done()
		listener, err := net.Listen("tcp", ":"+conf.Server.GrpcPort)
		failOnError(logger, err, "Failed to listen port")

		err = grpcServer.Serve(listener)
		failOnError(logger, err, "Failed to start GRPC server")
	}()

	<-mainShutCh
	logger.Infof("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	failOnError(logger, err, "failed to shut down HTTP server")

	grpcServer.GracefulStop()
	wg.Wait()
	close(shutDownCh)
	logger.Infof("Program terminated gracefully.")
}
