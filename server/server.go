package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type ChatServer struct {
	msg      map[string][]string
	users    []string
	shutdown chan bool
}

func main() {
	s := new(ChatServer)

	RunServer(s)

}

func (c *ChatServer) Register(username string, reply *string) error {
	*reply = "Welcome\n"
	*reply += "List of users online:\n"
	c.users = append(c.users, username)
	for _, value := range c.users {
		*reply += value + "\n"
	}

	for k := range c.msg {
		c.msg[k] = append(c.msg[k], username+" has joined.")
	}

	log.Printf("%s has joined the chat.\n", username)
	return nil
}

func RunServer(c *ChatServer) {
	rpc.Register(c)
	rpc.HandleHTTP()

	log.Printf("Listening on port :1234...\n")

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalln("error:", err)
	}

	http.Serve(l, nil)
}
