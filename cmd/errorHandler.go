package cmd

import (
	"fmt"
)

func FailOnError(err error, msg string) int {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		return 400
	}

	return 200
}
