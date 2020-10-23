package utils

import (
	"bufio"
	"os"
	"strings"
)

// Function that return 1 if you type y/Y/yes or 0 otherwise
func AskForConfirmation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	switch strings.ToLower(response) {
	case "y\n", "yes\n", "\n":
		return true, nil
	}
	return false, nil
}