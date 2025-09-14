package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gituhub.com/anton-jj/gator/internal/api"
	"gituhub.com/anton-jj/gator/internal/config"
	"gituhub.com/anton-jj/gator/internal/database"
	"os"
	"time"
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
		fmt.Errorf("failed to create/connect to db")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	appState.db = dbQueries

	cliCommands := commands{}
	cliCommands.register("login", handlerLogin)
	cliCommands.register("register", handlerRegister)
	cliCommands.register("reset", handlerReset)
	cliCommands.register("users", handlerUsers)
	cliCommands.register("agg", handleAgg)
	cliCommands.register("addfeed", handleAddFeed)
	cliCommands.register("feeds", handleFeeds)

	if len(os.Args) < 2 {
		fmt.Println("error give more arguments")
		os.Exit(1)
	}
	name := os.Args[1]
	args := []string{}
	if len(os.Args) == 2 {
		args = os.Args[1:]
	} else {

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
	return true
}

func handlerRegister(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("User %q already exists – please pick another name.\n", cmd.Args[0])
	}

	userParmas := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.db.CreateUser(context.Background(), userParmas)
	s.Cfg.SetUsername(cmd.Args[0])
	fmt.Println("user was created")
	fmt.Println("user name: ", userParmas.Name)
	fmt.Println("user uuid: ", userParmas.ID)
	fmt.Println("user created at: ", userParmas.CreatedAt)
	fmt.Println("user UpdatedAt: ", userParmas.UpdatedAt)

	return nil
}
func handlerLogin(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments\n")
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
func handlerReset(s *state, cmd command) error {
	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments\n")
	}
	err := s.db.ResetDatabase(context.Background())
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	fmt.Printf("The database was reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments\n")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Cfg.Current_Username {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	BASE_URL := "https://www.wagslane.dev/index.xml"

	agg, err := api.FetchFeed(context.Background(), BASE_URL)
	if err != nil {
		return err
	}

	fmt.Println(agg)
	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("to few args")
	}
	current, err := s.db.GetUser(context.Background(), s.Cfg.Current_Username)
	if err != nil {
		return err
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	addfeedParam := database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    current.ID,
	}

	s.db.CreateFeeds(context.Background(), addfeedParam)
	return nil
}

func handleFeeds(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments\n")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("feeds: \n")
	for _, feed := range feeds {
		fmt.Printf(" - name: %s\n - url: %s\n", feed.Name, feed.Url)
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf(" - created by: %s\n", user.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {

	if !checkArgLen(cmd) {
		return fmt.Errorf("pls send right arguments\n")
	}

	feed, err := s.db.FindFeedByUrl(context.Background(), cmd.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("%s ")

}
