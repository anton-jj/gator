package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gituhub.com/anton-jj/gator/internal/config"
	"gituhub.com/anton-jj/gator/internal/database"
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
	db  *database.Queries
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
	db, err := sql.Open("postgres", appState.Cfg.DB_URL)
	if err != nil {
		fmt.Println("failed to create/connect to db")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	appState.db = dbQueries

	cliCommands := commands{}
	cliCommands.register("login", handlerLogin)
	cliCommands.register("register", handlerRegister)

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

func checkArgLen(cmd command) bool {
	if len(cmd.Args) == 0 {
		return false
	}
	if len(cmd.Args) > 1 {
		return false
	}
	return true
}

func handlerRegister(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		errors.New("pls send right arguments")
	}

	fmt.Println(cmd.Args[0])
	userParmas := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.db.CreateUser(context.Background(), userParmas)
	fmt.Println("user was created")
	fmt.Println("user name: ", userParmas.Name)
	fmt.Println("user uuid: ", userParmas.ID)
	fmt.Println("user created at: ", userParmas.CreatedAt)
	fmt.Println("user UpdatedAt: ", userParmas.UpdatedAt)

	return nil
}
func handlerLogin(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		fmt.Printf("pls send right arguments\n")
	}

	user, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("user not found pls try again\n")
		return err
	}
	s.Cfg.SetUsername(user.Name)
	fmt.Printf("%s has logged in\n", user.Name)
	return nil

}
