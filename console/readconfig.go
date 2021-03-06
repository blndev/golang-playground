package main

// source https://github.com/BurntSushi/toml/tree/master/_examples

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title   string
	DB      database `toml:"database"`
	Servers map[string]server
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type server struct {
	IP string
	DC string
}

func main() {
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Title: %s\n", config.Title)

	fmt.Printf("Database: %s %v (Max conn. %d), Enabled? %v\n",
		config.DB.Server, config.DB.Ports, config.DB.ConnMax,
		config.DB.Enabled)
	for serverName, server := range config.Servers {
		fmt.Printf("Server: %s (%s, %s)\n", serverName, server.IP, server.DC)
	}
}
