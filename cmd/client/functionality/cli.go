package functionality

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/joho/godotenv"
)

var commands = make(map[string]handlerFunc)

func MainCLI() {
	initCommands()
	reader := bufio.NewReader(os.Stdin)
	newConf := &c.UserConfig{}
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("No login session found, please enter credentials.")
		newConf = initLogin(reader)
	} else {
		newConf.User = os.Getenv("CLI_USER")
		newConf.Token = os.Getenv("CLI_TOKEN")
		login(newConf)
	}

	for {
		fmt.Print("\nEnter a valid command.\n")
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
		}
		args := strings.Split(input, " ")
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		if args[0] == "login" {
			cmd := commands["login"]
			if len(args) > 1 {
				newConf.Args = args[1:]
			}
			success := cmd(newConf)
			newConf.Args = []string{}
			if !success {
				fmt.Println("Something went wrong.")
			}

		} else if cmd, ok := commands[args[0]]; ok {
			if len(args) > 1 {
				newConf.Args = args[1:]
			}
			success, err := authorizedMiddleware(newConf, cmd)
			newConf.Args = []string{} //Clear old args
			if err != nil {
				fmt.Println(err)
			}
			if !success {
				fmt.Println("Something went wrong.")
			}
		} else {
			fmt.Println("Currently unknown command.")
		}
	}
}

func initCommands() {
	commands["list"] = listcontents
	commands["login"] = login
}

func initLogin(reader *bufio.Reader) *c.UserConfig {

	for {
		fmt.Print("\nPlease enter login credentials or 'quit': ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
		}
		args := strings.Split(input, " ")
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		if args[0] == "quit" {
			os.Exit(1)
		}
		if len(args) > 1 {
			newConf := &c.UserConfig{Args: args}
			if login(newConf) {
				fmt.Println("Login successful.")
				return newConf
			} else {
				fmt.Println("Login unsuccesful. Please try again or 'quit'.")
			}
		} else {
			fmt.Println("Currently unknown command, please login.")
		}
	}
}
