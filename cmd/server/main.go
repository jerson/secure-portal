package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"secure-portal/modules/db"
	"secure-portal/services"

	"github.com/soheilhy/cmux"
)

func init() {
	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}
}
func migrate(ctx context.Context) {
	_, err := db.Setup(ctx)
	if err != nil {
		panic(err)
	}

}

func main() {

	ctx := context.NewContextSingle("main")
	defer ctx.Close()

	migrate(ctx)

	port := fmt.Sprintf(":%d", config.Vars.Server.Port)
	fmt.Println("running: ", port)

	rpcPort := fmt.Sprintf(":%d", config.Vars.Server.RPCPort)
	fmt.Println("running RPC: ", rpcPort)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	listenerRPC, err := net.Listen("tcp", rpcPort)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	services.RegisterServicesServer(server, &services.Server{})

	mux := cmux.New(listener)
	grpcListener := mux.Match(cmux.HTTP2())
	httpListener := mux.Match(cmux.HTTP1())

	group := new(errgroup.Group)
	group.Go(func() error { return grpcServe(server, listenerRPC) })
	group.Go(func() error { return grpcServe(server, grpcListener) })
	group.Go(func() error { return httpServe(httpListener) })
	group.Go(func() error { return mux.Serve() })

	err = group.Wait()
	if err != nil {
		panic(err)
	}
}

func httpServe(listen net.Listener) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("hi"))
	})

	httpServer := &http.Server{Handler: mux}

	return httpServer.Serve(listen)
}

func grpcServe(server *grpc.Server, listen net.Listener) error {
	return server.Serve(listen)
}
