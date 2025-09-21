package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gituhub.com/anton-jj/gator/internal/api"
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
	cfg, err := config.Read()
	if err != nil {
		return
	}
	appState := state{Cfg: &cfg}
	db, err := sql.Open("postgres", appState.Cfg.DB_URL)
	if err != nil {
		fmt.Errorf(err.Error())
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
	cliCommands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cliCommands.register("feeds", handleFeeds)
	cliCommands.register("follow", middlewareLoggedIn(handlerFollow))
	cliCommands.register("following", middlewareLoggedIn(handlerFollowing))
	cliCommands.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cliCommands.register("browse", middlewareLoggedIn(handleBrowse))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.Cfg.Current_Username)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
func checkArgLen(cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("pls send in an argument")
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if err := checkArgLen(cmd); err != nil {
		return err
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return err
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

	if err := checkArgLen(cmd); err != nil {
		return err
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
	if err := checkArgLen(cmd); err != nil {
		return err
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
	if err := checkArgLen(cmd); err != nil {
		return err
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
	time_between_req, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time_between_req)
	fmt.Printf("Collecting feeds every %v\n", time_between_req)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handleAddFeed(s *state, cmd command, user database.User) error {
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

	feed, err := s.db.CreateFeeds(context.Background(), addfeedParam)
	if err != nil {
		return err
	}
	newFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    current.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := s.db.CreateFeedFollow(context.Background(), newFeedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("feed: %s is now followed by %s\n", res.FeedName, res.UserName)
	return nil
}

func handleFeeds(s *state, cmd command) error {

	if err := checkArgLen(cmd); err != nil {
		return err
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if err := checkArgLen(cmd); err != nil {
		return err
	}
	feed, err := s.db.FindFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("failed at find feed by url?")
		return err
	}
	newFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res, err := s.db.CreateFeedFollow(context.Background(), newFeedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("feed: %s is now followed by %s\n", res.FeedName, res.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if err := checkArgLen(cmd); err != nil {
		return err
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println(err.Error())

		return err
	}

	fmt.Printf("%s is following: \n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("-%s\n", feed.FeedName)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if err := checkArgLen(cmd); err != nil {
		return err
	}
	fmt.Printf(cmd.Args[0])
	feed, err := s.db.FindFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	s.db.RemoveFeedFollowsRecord(context.Background(), database.RemoveFeedFollowsRecordParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	fmt.Printf("%s stopped following %s\n", user.Name, feed.Name)
	return nil

}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	res, err := api.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	id := uuid.NullUUID{
		UUID:  nextFeed.ID,
		Valid: true,
	}
	fmt.Println(id)
	for _, f := range res.Channel.Item {
		desc := sql.NullString{
			String: f.Description,
			Valid:  true,
		}
		fmt.Println(desc)
		date, err := time.Parse("YYYY-MM-DD", f.Description)
		if err != nil {
			return err
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       f.Title,
			Description: desc,
			PublishedAt: sql.NullTime{
				Time:  date,
				Valid: true,
			},
			Url:    f.Link,
			FeedID: id,
		}
		post, err := s.db.CreatePost(context.Background(), params)
		if err != nil {
			return err
		}
		fmt.Println(post)
	}
	return nil
}

func handleBrowse(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println(limit)

	if limit == 0 {
		limit = 2
	}

	fmt.Println(limit)
	for _, feed := range feeds {
		fmt.Println(feed.ID)
		params := database.GetPostsParams{
			FeedID: uuid.NullUUID{
				UUID:  feed.ID,
				Valid: true,
			},
			Limit: int32(limit),
		}
		fmt.Println(params)
		posts, err := s.db.GetPosts(context.Background(), params)
		if err != nil {
			return nil
		}
		fmt.Println(posts)
	}

	return nil
}
