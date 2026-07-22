package functionality

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	c "github.com/finfreezer/homeserver/internal/config"
	"github.com/joho/godotenv"
)

type Command struct {
	fun           handlerFunc
	requiresLogin bool
	name          string
	desc          string
}

var commands = make(map[string]Command)

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
		if cmd, ok := commands[args[0]]; ok {
			if len(args) > 1 {
				newConf.Args = args[1:]
			}
			success, err := authorizedMiddleware(newConf, cmd.fun, cmd.requiresLogin)
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
	commands["list"] = Command{
		fun:           listContents,
		requiresLogin: true,
		name:          "list",
		desc: `Lists all the files and subdirectories. Additional arguments 
		can be given to change the 'root' of the displayed tree. Starts from where the
		server assets are by default. list (optional){path/to/directory}`,
	}
	commands["login"] = Command{
		fun:           login,
		requiresLogin: false,
		name:          "login",
		desc: `Authorizes current user to allow use of commands. If the user has an
		existing log-in session, no arguments needed. 
		Otherwise use 'login {username} {password}`,
	}
	commands["stream"] = Command{
		fun:           streamContent,
		requiresLogin: true,
		name:          "stream",
		desc: `Opens a data stream of the desired content in the default browser.
		Current use 'stream {path/to/file}`,
	}
	commands["quit"] = Command{
		fun:           quitClient,
		requiresLogin: false,
		name:          "quit",
		desc:          "Quits the client.",
	}
	commands["help"] = Command{
		fun:           listCommands,
		requiresLogin: false,
		name:          "help",
		desc:          "Prints out the description for each available command.",
	}
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

func quitClient(cfg *c.UserConfig) bool {
	fmt.Printf("Goodbye, %s\n", cfg.User)
	os.Exit(0)
	return true
}

func listCommands(cfg *c.UserConfig) bool {
	for _, cmd := range commands {
		fmt.Printf("Description of '%s': %s\n\n", cmd.name, cmd.desc)
	}
	return true
}
