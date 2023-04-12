package cmd

import (
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/DragFAQ/uuid-generator/config"
	"github.com/DragFAQ/uuid-generator/generator"
	"github.com/DragFAQ/uuid-generator/handler"
	log "github.com/DragFAQ/uuid-generator/logger"
	pb "github.com/DragFAQ/uuid-generator/proto"
)

var (
	currentHash generator.Hash
	hashLock    sync.RWMutex
)

func failOnError(logger log.Logger, err error, msg string) {
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Panicf("%s: %s", msg, err)
	}
}

func setUpSignalHandler(ctx context.Context, wg *sync.WaitGroup, logger log.Logger, httpServer *http.Server, grpcServer *grpc.Server, cancel context.CancelFunc) {
	wg.Add(1)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-stop
		logger.Infof("shutting down (%v)", sig)

		err := httpServer.Shutdown(ctx)
		failOnError(logger, err, "failed to shut down HTTP server")
		grpcServer.GracefulStop()
		cancel()
		wg.Done()
	}()
}

func startHTTPServer(wg *sync.WaitGroup, port string, logger log.Logger) *http.Server {
	httpHandler := handler.NewHttpHandler(&currentHash, &hashLock, logger)
	http.HandleFunc("/", httpHandler.GetCurrentHash)

	srv := &http.Server{Addr: ":" + port}
	go func() {
		defer wg.Done()

		err := srv.ListenAndServe()
		failOnError(logger, err, "Failed to start HTTP server")
	}()

	return srv
}

func startGRPCServer(wg *sync.WaitGroup, port string, logger log.Logger) *grpc.Server {
	srv := grpc.NewServer()

	grpcHandler := handler.NewGrpcHandler(&currentHash, &hashLock, logger)
	pb.RegisterHashServiceServer(srv, grpcHandler)

	listener, err := net.Listen("tcp", ":"+port)
	failOnError(logger, err, "Failed to listen port")

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = srv.Serve(listener)
		failOnError(logger, err, "Failed to start GRPC server")
	}()

	return srv
}

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

			ctx, cancel := context.WithCancel(context.Background())

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				generator.GenerateHash(&currentHash, &hashLock, logger, conf.Settings.HashTTLSeconds, ctx)
			}()

			httpServer := startHTTPServer(wg, conf.Server.HttpPort, logger)

			grpcServer := startGRPCServer(wg, conf.Server.GrpcPort, logger)

			setUpSignalHandler(ctx, wg, logger, httpServer, grpcServer, cancel)

			wg.Wait()
			logger.Infof("program terminated gracefully")
		},
	}

	return cmd
}
