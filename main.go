package main

import (
	"fmt"
	"gituhub.com/anton-jj/gator/internal/config"
)

func main() {
	gatorConf, err := config.Read()
	if err != nil {
		return
	}
	gatorConf.SetUsername("anton")
	gatorConf.Write(gatorConf)
	fmt.Printf("%s\n%s\n", gatorConf.Current_Username, gatorConf.DB_URL)

}
