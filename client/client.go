package main

import (
	"flag"
	"log"
	"net/rpc"
	"os"
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
	if len(os.Args) < 3 {
		log.Fatalln("Usage go run client.go <name> <message>")
		os.Exit(1)
	}

	client.Registre()
	// err = client.Call("ChatServer.List", args, &reply2)
	// if err != nil {
	// 	log.Fatalln("Error: ", err)
	// }
	// fmt.Printf("%v\nlist: %v\n", reply, reply2)
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
