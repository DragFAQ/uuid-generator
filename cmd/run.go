package cmd

import (
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/DragFAQ/uuid-generator/config"
	"github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
)

const shutdownTimeout = 10 * time.Second

var (
	currentHash generator.Hash
	hashLock    sync.RWMutex
	wg          sync.WaitGroup
)

func failOnError(logger log.Logger, err error, msg string) {
	if err != nil && err != http.ErrServerClosed {
		logger.Panicf("%s: %s", msg, err)
	}
}

func setUpSignalHandler(logger log.Logger, httpServer *http.Server, grpcServer *grpc.Server) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		sig := <-signalCh
		logger.Infof("shutting down (%v)", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		failOnError(logger, err, "failed to shut down HTTP server")

		grpcServer.GracefulStop()
		wg.Wait()

		close(signalCh)
		logger.Infof("program terminated gracefully")
		os.Exit(0)
	}()
}

// Run main cmd action
func Run() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run uuid-generator",
		Long:  `Run uuid-generator as HTTP and GRPC server.`,
		Run: func(cmd *cobra.Command, args []string) {
			var conf = config.NewConfig()

			var logger, err = log.NewZapLogger(conf.Logger)
			failOnError(logger, err, "new logger")

			logger.Infof("starting service on %v", time.Now())

			// Start generating the initial hash
			currentHash = generator.Hash{
				Value:          uuid.New().String(),
				GenerationTime: time.Now(),
			}

			generateCh := make(chan os.Signal, 1)
			signal.Notify(generateCh, syscall.SIGINT, syscall.SIGTERM)

			wg.Add(1)
			go generator.GenerateHash(&currentHash, &hashLock, logger, conf.Settings.HashTTLSeconds, generateCh, &wg)

			// Define the HTTP API routes
			httpHandler := handler.NewHttpHandler(&currentHash, &hashLock, logger)
			http.HandleFunc("/", httpHandler.GetCurrentHash)

			httpServer := &http.Server{Addr: ":" + conf.Server.HttpPort}
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := httpServer.ListenAndServe()
				failOnError(logger, err, "Failed to start HTTP server")
			}()

			// Define the gRPC server
			grpcServer := grpc.NewServer()

			grpcHandler := handler.NewGrpcHandler(&currentHash, &hashLock, logger)
			pb.RegisterHashServiceServer(grpcServer, grpcHandler)

			listener, err := net.Listen("tcp", ":"+conf.Server.GrpcPort)
			failOnError(logger, err, "Failed to listen port")

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = grpcServer.Serve(listener)
				failOnError(logger, err, "Failed to start GRPC server")
			}()

			setUpSignalHandler(logger, httpServer, grpcServer)
			select {}
		},
	}

	return cmd
}
