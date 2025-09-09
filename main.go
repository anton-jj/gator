package main

import (
	"errors"
	"fmt"
	"os"

	"gituhub.com/anton-jj/gator/internal/config"
)

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.Handlers[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("something went wrong")
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*state, command) error)
	}
	c.Handlers[name] = f
}

type state struct {
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

func main() {
	// gatorConf, err := config.Read()
	cfg, err := config.Read()
	if err != nil {
		return
	}
	appState := state{Cfg: &cfg}
	cliCommands := commands{}
	cliCommands.register("login", handlerlogin)

	fmt.Println(os.Args[2:])
	if len(os.Args) < 2 {
		fmt.Println("error give more arguments")
		os.Exit(1)
	}
	name := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	cmd := command{Name: name, Args: args}
	if err := cliCommands.run(&appState, cmd); err != nil {
		os.Exit(1)
	}

}

func handlerlogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		errors.New("hello send some arguments")
	}
	if len(cmd.Args) > 1 {
		errors.New("pls my friend one argument is enough")
	}
	s.Cfg.SetUsername(cmd.Args[0])
	fmt.Printf("the username has ben set to %s\n", cmd.Args[0])
	return nil

}
