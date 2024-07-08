package utils

import (
	"fmt"
	"log"
)

// the bool is to determine if the user wants to exit that menu
func GetUserResponse(prompt string, formatVariables ...any) (string, bool) {
	fmt.Printf(prompt+"\n", formatVariables...)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal("Failed to parse user input" + err.Error())
	}
	if response == "0" {
		return response, true
	}
	return response, false
}
