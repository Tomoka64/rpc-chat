package main

import (
	"flag"
	"log"
	"net/rpc"
)

type Chatclient struct {
	Username string
	Client   *rpc.Client
}

func main() {
	client, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}

	client.Registre()

}

func newClient() (*Chatclient, error) {
	var c *Chatclient = &Chatclient{}
	flag.StringVar(&c.Username, "user", "tomoka", "Your username")
	flag.Parse()
	return c, nil

}

func (c *Chatclient) getConn() *rpc.Client {
	var err error

	if c.Client == nil {
		c.Client, err = rpc.DialHTTP("tcp", ":1234")
		if err != nil {
			log.Panicf("unable to make connection: %q", err)
		}
	}
	return c.Client
}

func (c *Chatclient) Registre() {
	var reply string
	c.Client = c.getConn()

	err := c.Client.Call("ChatServer.Register", c.Username, &reply)
	if err != nil {
		log.Printf("Error registering user: %q", err)
	} else {
		log.Printf("Reply: %s", reply)
	}
}
