// File: https://github.com/data-engineering-helpers/dppctl/blob/main/tests/check-dppctl.go

package main

import (
  //"fmt"
  "log"
  "github.com/data-engineering-helpers/dppctl/service"
)

func main() {
  // Set properties of the predefined Logger, including
  // the log entry prefix and a flag to disable printing
  // the time, source file, and line number.
  log.SetPrefix("[test-dppctl] ")
  log.SetFlags(0)

  // AWS S3
  file_list, err := service.S3List("baldwins")
  if err != nil {
    log.Print(err)
  }

  for _, file_metadata := range file_list {
    log.Println(file_metadata)
  }
}


