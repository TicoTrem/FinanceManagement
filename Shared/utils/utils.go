package utils

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
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

func GetUserResponseInt(prompt string, lowest int, highest int, formatVariables ...any) (parsedInt int, createNew bool, exit bool) {
	for {
		response, exit := GetUserResponse(prompt, formatVariables...)
		if exit {
			return -1, false, true
		}
		if response == "c" {
			return -1, true, false
		}
		pInt, err := strconv.Atoi(response)
		if err != nil || pInt < 0 || pInt < lowest || pInt > highest {
			fmt.Println("Invalid Input")
			continue
		}
		return pInt, false, exit
	}
}

// prints out records of any type of struct
func SelectRecordOrCreate[T any](records []T, createNewFunc func()) *T {
	// get the value of the first struct

	value := reflect.ValueOf(records[0])

	// print all of the records
	for i := 0; i < len(records); i++ {

		if value.Kind() != reflect.Struct {
			log.Fatal("This interface was not a struct")
			return nil
		}

		// get the struct type of the first struct
		structType := value.Type()

		var structString string
		for i := 0; i < value.NumField(); i++ {
			fieldName := structType.Field(i).Name
			fieldValue := value.Field(i).Interface()
			structString += fmt.Sprintf("%v: %v\t", fieldName, fieldValue)
		}
		fmt.Println(structString)
	}

	pInt, createNew, exit := GetUserResponseInt("Enter the number of the record you would like to edit, or 'c' to create a new one", 1, len(records))
	if exit {
		return nil
	}
	if createNew {
		if createNewFunc == nil {
			log.Fatal("The create new function was nil and not callable")
		}
		createNewFunc()
		return nil
	}

	// return the selected option to where we know everything about the passed in objects
	return &records[pInt-1]
}

// This method will take a prompt string that will be displayed to the user, as well as a slice of strings
// that correspond to the positions of the method they will call in methodsToCall. If the user answers 'c'
// the createNewFunc will be called to handle creating a new record. If you do not need the option to create
// a new record, just pass in 'nil' to that parameter. functions are just pointers so this will not cause errors
func PromptAndHandle(prompt string, options []string, methodsToCall []func(), formatVariables ...any) {
	for {
		for i := 0; i < len(options); i++ {
			prompt += fmt.Sprintf("%v)\t%v\n", i+1, options[i])
		}
		fmt.Printf(prompt)

		reader := bufio.NewReader(os.Stdin)
		userResponse, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err)
			return
		}
		userResponse = strings.TrimSpace(userResponse)
		if userResponse == "exit" {
			return
		}
		pInt, err := strconv.Atoi(userResponse)
		if err != nil || pInt < 1 || pInt > len(options) {
			fmt.Println("Invalid Input")
			continue
		}

		// call the corresponding method
		methodsToCall[pInt-1]()
		break
	}
}

//func CreateNewOrInt(prompt string, minimum int, maximum int, formatVariables ...any) (response int, createNew bool, exit bool) {
//	for {
//		response, exit := GetUserResponse(prompt, formatVariables...)
//		if exit {
//			return -1, false, true
//		}
//		lowercase := strings.ToLower(response)
//		if lowercase == "c" {
//			return -1, true, exit
//		} else {
//			parsedInt, err := strconv.Atoi(response)
//			if err != nil || parsedInt < minimum || parsedInt > maximum {
//				fmt.Println("Invalid Input")
//				continue
//			}
//			return parsedInt, false, exit
//		}
//	}
//}
//

// returns a negative sign in front of the $ if the float is negative
func GetMoneyString(money float32) string {
	var estimatedMoneyString string

	if money < 0 {
		absEstimatedSpendingMoney := math.Abs(float64(money))
		estimatedMoneyString = fmt.Sprintf("-$%.2f", absEstimatedSpendingMoney)
	} else {
		estimatedMoneyString = fmt.Sprintf("$%.2f", money)
	}
	return estimatedMoneyString
}
