package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework10/internal/adapters/repo"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcPort "homework10/internal/ports/grpc"
)

const (
	gPort = ":50054"
	hPort = ":9000"
)

func main() {
	lis, err := net.Listen("tcp", gPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	adApp := app.NewApp(repo.New(), repo.New())

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(grpcPort.InterceptorLogger()),
		recovery.UnaryServerInterceptor([]recovery.Option{
			recovery.WithRecoveryHandler(grpcPort.PanicInterceptor),
		}...),
	))
	grpcService := grpcPort.NewService(adApp)
	grpcPort.RegisterAdServiceServer(grpcServer, grpcService)

	httpServer := httpgin.NewHTTPServer(hPort, adApp)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		log.Printf("starting grpc cmd, listening on %s\n", gPort)
		defer log.Printf("close grpc cmd listening on %s\n", gPort)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc cmd can't listen and serve requests: %w", err)
		}
	})

	eg.Go(func() error {
		log.Printf("starting http cmd, listening on %s\n", httpServer.Addr)
		defer log.Printf("close http cmd listening on %s\n", httpServer.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http cmd listening on %s: %s", httpServer.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http cmd can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}
