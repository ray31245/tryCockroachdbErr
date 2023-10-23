package crossservice

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/cockroachdb/errors/grpc/middleware"
	"github.com/hydrogen18/memlistener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Client EchoerClient
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	srv := &EchoServer{}

	lis := memlistener.NewMemoryListener()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryServerInterceptor))
	RegisterEchoerServer(grpcServer, srv)

	go grpcServer.Serve(lis)

	dialOpts := []grpc.DialOption{
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial("", "")
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor),
	}

	clientConn, err := grpc.Dial("", dialOpts...)
	if err != nil {
		panic(err)
	}

	Client = NewEchoerClient(clientConn)

	code := m.Run()

	grpcServer.Stop()
	clientConn.Close()

	os.Exit(code)
}
