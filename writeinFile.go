package main

import (
  "bufio"
  "os"
  "fmt"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func WriteInFile(s string, fileName string) {
  f, err := os.Create(fileName)
  check(err)
  // It's idiomatic to defer a `Close` immediately
  // after opening a file.
  defer f.Close()
  w := bufio.NewWriter(f)
  n4, errw := w.WriteString(string(s))
  check(errw)
  fmt.Printf("wrote %d bytes\n", n4)
  // Use `Flush` to ensure all buffered operations have
  // been applied to the underlying writer.
  w.Flush()
}
