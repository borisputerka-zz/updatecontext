package utils

import (
	"bufio"
	"os"
	"strings"
)

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

func StringInSlice(a string, list []string) bool {
	for _, item := range list {
		if item == a {
			return true
		}
	}
	return false
}
