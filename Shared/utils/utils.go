package utils

import (
	"fmt"
	"log"
	"strconv"
)

// the bool is to determine if the user wants to exit that menu
func GetUserResponse(prompt string, formatVariables ...any) (response string, exit bool) {
	fmt.Printf(prompt+"\n", formatVariables...)
	var userResponse string
	_, err := fmt.Scanln(&userResponse)
	if err != nil {
		log.Fatal("Failed to parse user input" + err.Error())
	}
	if userResponse == "0" {
		return response, true
	}
	return response, false
}

// the extra bool on top of GetUserResponse that this returns is meant to tell the
// caller if the parsing was successful
func GetUserResponseFloat(prompt string, formatVariables ...any) (parsedFloat float32, exit bool) {
	for {
		response, exit := GetUserResponse(prompt, formatVariables...)
		pFloat, err := strconv.ParseFloat(response, 32)
		if err != nil {
			fmt.Println("Invalid Input")
			continue
		}
		return float32(pFloat), exit
	}
}

func GetUserResponseInt(prompt string, formatVariables ...any) (parsedInt int, exit bool) {
	for {
		response, exit := GetUserResponse(prompt, formatVariables...)
		pInt, err := strconv.Atoi(response)
		if err != nil {
			fmt.Println("Invalid Input")
			continue
		}
		return pInt, exit
	}
}
