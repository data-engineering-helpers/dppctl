// File: https://github.com/data-engineering-helpers/dppctl/blob/main/go/tests/check-dppctl.go

package main

import (
  //"fmt"
  "log"
  "github.com/data-engineering-helpers/dppctl"
)

func main() {
  // Set properties of the predefined Logger, including
  // the log entry prefix and a flag to disable printing
  // the time, source file, and line number.
  log.SetPrefix("[test-dppctl] ")
  log.SetFlags(0)

  // Test with an empty name, which should trigger an error
  _, err0 := dppctl.Hello("")
  // If an error was returned, print it to the console
  if err0 != nil {
      log.Print(err0)
  }

  //
  name := "Test of Data Processing Pipeline (DPP) CLI utility"
  message, _ := dppctl.Hello(name)
  log.Print(message)
}


