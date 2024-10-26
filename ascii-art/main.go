package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	receive := os.Args
	text_to_read := receive[1]
	myError := errors.New("Error!!!")
	if text_to_read == "" {
		return
	}
	if len(receive) != 2 {
		fmt.Println(myError)
		return
	}
	if text_to_read == "\\n" {
		fmt.Println()
		return
	}
	text_in_array := strings.Split(text_to_read, `\n`)
	var line []string
	readFile, err := os.Open("standard.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}
	var line_row []string
	print_text(text_in_array, line_row, line)
}

func print_text(text_in_array []string, line_row []string, line []string) {
	for array_index := 0; array_index < len(text_in_array); array_index++ {
		word := text_in_array[array_index]
		line_row = art_char(word, text_in_array, array_index, line_row, line)
		if len(line_row) == 0 {
		} else {
			for i := 0; i < 9; i++ {
				for k := i; k < len(line_row); k += 9 {
					fmt.Print(line_row[k])
				}
				if i > 0 && i < 8 {
					fmt.Println()
				}
			}
		}
		fmt.Println()
		line_row = nil
	}
}

func art_char(word string, text_in_array []string, array_index int, line_row []string, line []string) []string {
	for symbol_index := 0; symbol_index < len(text_in_array[array_index]); symbol_index++ {
		symbol := int(word[symbol_index])
		if symbol >= 32 && symbol < 127 {
			start_read := (symbol-32)*8 + (symbol - 32)
			end_read := start_read + 8
			for lines := start_read; lines <= end_read; lines++ {
				line_row = append(line_row, line[lines])
			}
		}
	}
	return line_row
}
