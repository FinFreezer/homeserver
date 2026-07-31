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
	commands["listall"] = Command{
		fun:           listAllContents,
		requiresLogin: true,
		name:          "listall",
		desc: `Lists all the files and subdirectories starting from the working directory,
		up to a depth of 99 branches. With the -dirOnly flag, skips displaying files.
		Additional arguments can be given to change the 'root' of the displayed tree.
		Starts from current working directory by default. 'listall {-dironly} {pathToRoot}'
		Alternatively you can set the depth of the listing with 'listall {-d} {num} {pathToRoot}'.`,
	}
	commands["list"] = Command{
		fun:           listContents,
		requiresLogin: true,
		name:          "list",
		desc: `Lists all the files and subdirectories of the current working directory.
		Additional -dironly flag displays only subdirectories.
		'list {-dironly}.`,
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
		desc: `Opens a data stream of the desired content in VLC.
		Current use 'stream {path/to/file.png}', opening a browser stream by default. Optional flags include '-f' for playing single files,
		'-b' to open the target in a browser, and '-a' to receive a playlist of the target folder, recommended
		for usage with media players. To request a response with a link to a playlist instead of being redirected
		directly to a media player, use the '-po' (playlist only) flag. 'stream {-a|-b|-f|-po} {path/to/file.png}`,
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
	commands["cd"] = Command{
		fun:           changeDirectory,
		requiresLogin: true,
		name:          "cd",
		desc: `Changes the root folder for commands within the server. 'cd ..' 
		to go up a directory. cd {pathToSubfolder/a/b...} to go down a directory.`,
	}
	commands["update"] = Command{
		fun:           remoteUpdate,
		requiresLogin: true,
		name:          "update",
		desc: `Runs a remote script that updates the server to the newest version,
		and restarts it on the host machine. Note that you must modify the script on the
		host machine to include the right log-in credentials.`,
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
