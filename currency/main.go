package currency

import (
	"fmt"
	protos "github.com/MykhailoKondrat/go-microservices/currency/protos"
	"github.com/MykhailoKondrat/go-microservices/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	log := hclog.Default()
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)
	protos.RegisterCurrencyServer(gs, cs)
	l, err := net.Listen("tcp", ":9092")
	reflection.Register(gs)
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	err = gs.Serve(l)
	fmt.Println(err)
	if err != nil {
		fmt.Println()
	}

}
