package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/entity"
	"google.golang.org/grpc"
)

func main() {
	host := flag.String("host", "127.0.0.1", "grafana gRPC server host")
	port := flag.String("port", "10000", "grafana gRPC server port")
	token := flag.String("token", "", "grafana API token")
	flag.Parse()
	if token == nil || *token == "" {
		fmt.Println("Please provide a grafana API token")
		flag.PrintDefaults()
		return
	}

	addr := net.JoinHostPort(*host, *port)
	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
	}

	conn, err := grpc.DialContext(context.Background(), addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(tokenAuth{
			token: *token,
		}),
	)
	if err != nil {
		panic(err)
	}
	client := entity.NewEntityStoreClient(conn)

	fmt.Printf("listening on %s\n", config.Address)
	engine := sqle.NewDefault(sql.NewDatabaseProvider(NewDatabase(client)))
	engine.Analyzer.Debug = true
	engine.Analyzer.Verbose = true

	server, err := server.NewDefaultServer(config, engine)
	if err != nil {
		return
	}
	defer server.Close()

	server.Start()
}
