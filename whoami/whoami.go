package whoami

import (
	"fmt"
	"os/user"
)

func Main(argv []string) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", user.Name)
}
