package main

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

const portNumber = ":8080"

// Data is a struct for sending Ascii art and errors to html templates
type Data struct {
	AsciiArt string
	Error    error
}

func main() {
	fmt.Printf("Starting application on http://localhost%s/\n", portNumber)

	http.HandleFunc("/", Home)
	http.HandleFunc("/ascii-art", AsciiArt)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	log.Fatal(http.ListenAndServe(portNumber, nil))
}

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}
	if renderTemplate(w, "index.html", nil) != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
}

// AsciiArt is the printed Ascii art page handler
func AsciiArt(w http.ResponseWriter, r *http.Request) {
	data := Data{}
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}
	// Check user inputg
	userInput, err := checkInputCharacters(r.FormValue("user-input"))
	if err != nil {
		data.Error = err
		http.Redirect(w, r, "/", http.StatusBadRequest)
		renderTemplate(w, "index.html", data) // Render errors to html template
		return
	}
	banner := r.FormValue("banner")
	data.AsciiArt, err = makeAsciiArt(userInput, banner)
	if err != nil {
		data.Error = err
		http.Redirect(w, r, "/", http.StatusNotFound)
		renderTemplate(w, "index.html", data) // Render errors to html template
		return
	}
	if renderTemplate(w, "index.html", data) != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	} // Render Ascii art to html template
}

// renderTemplate renders templates using html/template
func renderTemplate(w http.ResponseWriter, templateName string, data any) error {
	parsedTemplate, err := template.ParseFiles("./templates/" + templateName)
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return err
	}
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template: ", err)
		return err
	}
	return nil
}

// checkInputCharacters checks user input for errors
func checkInputCharacters(userInput string) (userInputSlice []string, err error) {
	if userInput == "" {
		err = fmt.Errorf("Input cannot be empty. Please enter text")
		return []string{}, err
	}
	if len(userInput) > 1 && userInput[0] == '"' && userInput[len(userInput)-1] == '"' {
		userInput = userInput[1 : len(userInput)-1]
	}
	word := ""
	for _, char := range userInput {
		if char == 13 { // New line
			userInputSlice = append(userInputSlice, word)
			userInputSlice = append(userInputSlice, "")
			word = ""
			continue
		} else if char == 10 {
			continue
		} else if char < 32 || 126 < char { // Invalid character
			err = errors.New(fmt.Sprintf("%c is not a valid character. Please enter ASCII characters between 32-126\n", char))
			return []string{}, err
		} else {
			word += string(char)
		}
	}
	userInputSlice = append(userInputSlice, word) // Append last word
	return userInputSlice, nil
}

// makeAsciiArt converts user input text to Ascii art
func makeAsciiArt(userInput []string, fileName string) (output string, err error) {
	// Open file for [BANNER]
	file, err := os.Open("./fonts/" + fileName + ".txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read [BANNER] file lines to asciiCharacters []string{}
	scanner := bufio.NewScanner(file)
	asciiCharacters := []string{}
	for scanner.Scan() {
		if scanner.Text() != "" {
			asciiCharacters = append(asciiCharacters, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return printAscii(userInput, asciiCharacters), nil
}

// printAscii takse user input and returns Ascii art as a string
func printAscii(userInput []string, asciiCharacters []string) (output string) {
	for _, word := range userInput {
		if word == "" {
			output += "\n"
		} else {
			for i := 0; i < 8; i++ {
				for _, char := range word {
					output += asciiCharacters[(int(char)-32)*8+i]
				}
				if i != 7 {
					output += "\n"
				}
			}
		}
	}
	println(output)
	return output
}
