package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type N bool

type ChatServer struct {
	msg      map[string][]string
	users    []string
	shutdown chan bool
}

type Message struct {
	User   string
	Target string
	// ToAll  []string
	Msg string
}

func main() {
	s := new(ChatServer)
	s.msg = make(map[string][]string)
	s.shutdown = make(chan bool, 1)
	RunServer(s)
	<-s.shutdown
}

func (c *ChatServer) Register(username string, reply *string) error {
	*reply = "Welcome\n"
	*reply += "==========================\nList of users:\n"
	c.users = append(c.users, username)
	for _, value := range c.users {
		*reply += value + "\n"
	}
	*reply += "========================="

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

	go http.Serve(l, nil)
}

func (c *ChatServer) CheckMsgs(username string, reply *[]string) error {
	*reply = c.msg[username]
	c.msg[username] = nil
	return nil
}

func (c *ChatServer) Tell(msg Message, reply *N) error {

	if q, ok := c.msg[msg.Target]; ok {
		m := msg.User + " tells you " + msg.Msg
		c.msg[msg.Target] = append(q, m)
	} else {
		m := msg.Target + " does not exist"
		c.msg[msg.User] = append(q, m)
	}
	log.Println(msg.User + " said " + msg.Target + " '" + msg.Msg + "'")
	*reply = false

	return nil
}

func (c *ChatServer) TellAll(msg Message, reply *N) error {

	for k, v := range c.msg {
		m := msg.User + " says " + msg.Msg
		c.msg[k] = append(v, m)
	}
	log.Println(msg.User + " said all '" + msg.Msg + "'")
	*reply = true

	return nil
}
