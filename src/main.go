package main

// bjira ls -q sprint_name="MW Sprint 33" -q user="marcin.boczek"

import (
	"fmt"
	"log"
	"os"
	"strings"
	"runtime"
	"path/filepath"
)

type Config struct {
	URL   string
	email string
	token string 
}

var global_config Config = Config{ "null", "null", "null" }

func feed_config(path string, cfg Config) {
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	f_output := make([]byte, info.Size() + 1)
	file.Read(f_output)

	lines := strings.Split(string(f_output), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "JIRA_URL=") {
			global_config.URL = strings.TrimPrefix(line, "JIRA_URL=")
		}else if strings.HasPrefix(line, "JIRA_EMAIL=") {
			global_config.email = strings.TrimPrefix(line, "JIRA_EMAIL=")
		} else if strings.HasPrefix(line, "JIRA_API_TOKEN=") {
			global_config.token = strings.TrimPrefix(line, "JIRA_API_TOKEN=")
		}
	}
}

type Command struct {
	comm_str string ""
	description string ""
	comm_run func() error 
}

var command_array = []Command{
    {
		"ls", 
		"List all tasks based on query.", 
		ls_command,
	},
}
func help_command() {
	for _, comm := range command_array {
		fmt.Printf("\t%s - %s\n", comm.comm_str, comm.description)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bjira [command] [options]")
		help_command()
		return
	}

	if runtime.GOOS == "windows" {
		exe_path, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		cfg_path := filepath.Join(filepath.Dir(exe_path), "cfg")
		feed_config(cfg_path, global_config)
	} else if runtime.GOOS == "linux" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		cfg_path := filepath.Join(home, ".config", "bjira", "cfg")
		feed_config(cfg_path, global_config)
	}

	// Take input from user
	user_command := os.Args[1]
	for _, comm := range command_array {
		if comm.comm_str == user_command {
			err := comm.comm_run();
			if err == nil {
				return 
			} else {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("Usage: bjira [command] [options]")
	os.Exit(1)
}
