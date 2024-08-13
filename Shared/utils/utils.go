package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// the bool is to determine if the user wants to exit that menu
func GetUserResponse(prompt string, formatVariables ...any) (response string, exit bool) {
	fmt.Printf(prompt+"\n", formatVariables...)
	reader := bufio.NewReader(os.Stdin)
	userResponse, err := reader.ReadString('\n')
	userResponse = strings.TrimSpace(userResponse)

	if err != nil {
		log.Fatal("Error reading input:", err)
		return
	}
	if userResponse == "exit" {
		return userResponse, true
	}
	return userResponse, false
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
		if err != nil || pInt < 0 {
			fmt.Println("Invalid Input")
			continue
		}
		return pInt, exit
	}
}

func CreateNewOrInt(prompt string, minimum int, maximum int, formatVariables ...any) (response int, createNew bool, exit bool) {
	for {
		response, exit := GetUserResponse(prompt, formatVariables...)
		if exit {
			return -1, false, true
		}
		lowercase := strings.ToLower(response)
		if lowercase == "c" {
			return -1, true, exit
		} else {
			parsedInt, err := strconv.Atoi(response)
			if err != nil || parsedInt < minimum || parsedInt > maximum {
				fmt.Println("Invalid Input")
				continue
			}
			return parsedInt, false, exit
		}
	}
}
