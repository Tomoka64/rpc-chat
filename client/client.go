package main

import (
	"bufio"
	"flag"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Chatclient struct {
	Username string
	Client   *rpc.Client
}
type Message struct {
	User   string
	Target string
	// ToAll  []string
	Msg string
}

func main() {

	client, err := newClient()
	if err != nil {
		log.Fatalln(err)
	}

	client.Registre()
	go client.CheckMessages()

	loop(client)
}

func (c *Chatclient) CheckMessages() {
	var reply []string
	c.Client = c.getConn()

	for {
		err := c.Client.Call("ChatServer.CheckMsgs", c.Username, &reply)
		if err != nil {
			log.Fatalln("Chat has been shutdown.")

		}

		for i := range reply {
			log.Println(reply[i])
		}

		time.Sleep(time.Second)
	}
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

func loop(c *Chatclient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error: %q\n", err)
		}
		line = strings.TrimSpace(line)
		p := strings.Fields(line)
		if strings.HasPrefix(strings.Join(p, ""), "tell") == true {
			c.Tell(p)
		} else if strings.HasPrefix(strings.Join(p, ""), "toAll") == true {
			c.TellAll(p)
		} else {
			color.Red("you wanna tell someone specific? \n")
			color.Yellow("<usage> tell <username> <msg>\n")
			color.Red("you wanna tell everyone? <usage> toAll <msg>")
			color.Yellow("<usage> toAll <msg>\n")
		}
	}
}

type N bool

func (c *Chatclient) Tell(p []string) {
	var reply N
	c.Client = c.getConn()
	if len(p) > 2 {
		msg := strings.Join(p[2:], " ")
		message := Message{
			User:   c.Username,
			Target: p[1],
			Msg:    msg,
		}
		err := c.Client.Call("ChatServer.Tell", message, &reply)
		if err != nil {
			log.Printf("Error telling users something: %q", err)
		}

	}

}

func (c *Chatclient) TellAll(params []string) {
	var reply N
	c.Client = c.getConn()

	if len(params) > 1 {
		msg := strings.Join(params[1:], " ")
		message := Message{
			User: c.Username,
			Msg:  msg,
		}

		err := c.Client.Call("ChatServer.TellAll", message, &reply)
		if err != nil {
			log.Printf("Error: %q", err)
		}
	} else {
		log.Println("sent it to everyone")
	}

}
