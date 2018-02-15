package model

import "log"

type ChatServer struct {
	msg      map[string][]string
	users    []string
	shutdown chan bool
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
