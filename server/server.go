package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Tomoka64/gRCP/model"
)

func main() {
	s := new(model.ChatServer)

	RunServer(s)

}

func RunServer(c *model.ChatServer) {
	rpc.Register(c)
	rpc.HandleHTTP()

	log.Printf("Listening on port :1234...\n")

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalln("error:", err)
	}

	http.Serve(l, nil)
}
