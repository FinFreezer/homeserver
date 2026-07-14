package functionality

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	c "github.com/finfreezer/homeserver/internal/config"
)

func MainCLI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter a valid command.\n")
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
		}
		args := strings.Split(input, " ")
		for i, _ := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		if args[0] == "login" && len(args) == 3 {
			newConf := c.Read(args[1], args[2])
			if login(newConf) {
				fmt.Println("Login success.")
			} else {
				fmt.Println("Login unsuccesful.")
				break
			}
		} else {
			fmt.Println("Currently unknown command.")
		}
	}
}
