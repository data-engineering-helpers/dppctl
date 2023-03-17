// File: https://github.com/data-engineering-helpers/dppctl/blob/main/go/dppctl/dppctl.go

package dppctl

import (
  "fmt"
  "errors"
)

import "rsc.io/quote"

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
    // If no name was given, return an error with a message.
    if name == "" {
        return "", errors.New("empty name")
    }

    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message, nil
}

func main() {
  name := "Data Processing Pipeline (DPP) CLI utility"
  fmt.Println(Hello(name))
  fmt.Println(quote.Go())
}

