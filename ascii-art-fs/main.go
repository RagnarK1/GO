package main

import (
	"fmt"
	"os"
	"log"
	"strings"
)

func main(){
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments.")
		fmt.Println("Usage: go run . [STRING] [BANNER]\nEX: go run . something standard")
	} else {
		input := args[0]
		banner := args[1]
		if !validateInput(input){
			log.Fatalln("Non-ASCII character detected. Please, enter valid input.")
		}
		if !validateBanner(banner){
			log.Fatalln("This banner does not exist. Please, use 'standard', 'shadow', 'thinkertoy' .")
		}
		fmt.Println(makeAsciiArt(input, banner))
	}
}

func validateInput(s string) bool {
	check := true
	for _, symbol := range s {
		if symbol < 32 || symbol > 126 {
			check = false
		}
	}
	return check
}

func validateBanner(b string) bool {
	if b == "standard" || b == "shadow" || b == "thinkertoy" {
		return true
	}
	return false
}

func makeAsciiArt(input string, banner string) string {
	var inputArr []string
	inputArr = splitByNewline(input, inputArr)
	file, _ := os.ReadFile("fonts/" + banner + ".txt")
	file = []byte(strings.ReplaceAll(string(file), "\r", ""))
	asciiArr := strings.Split(string(file), "\n")
	var result []byte

	for r := 0; r < len(inputArr); r++ {
		if len(inputArr[r]) != 0 {
			for k := 0; k < 8; k++ {
				for i := 0; i < len(inputArr[r]); i++ {
					result = append(result, asciiArr[(int(inputArr[r][i])-31)*9-8+k]...)
				}
				result = append(result, byte(10))
			}
		} else {
			result = append(result, byte(10))
		}
	}
	return string(result[:len(result)-1])
}

func splitByNewline(input string, inputArr []string) []string {
	index := strings.Index(input, "\\n")
	if index == -1 {
		inputArr = append(inputArr, input)
		return inputArr
	} 
	left := input[:index]
	right := input[index+2:]
	inputArr = append(inputArr, left)
	return splitByNewline(right, inputArr)	
}