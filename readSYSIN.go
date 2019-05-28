// package main
//
// import (
//   "bufio"
//   "fmt"
//   "os"
//   "strings"
// )
//
// func main() {
//
//   reader := bufio.NewReader(os.Stdin)
//   fmt.Println("Simple Shell")
//   fmt.Println("---------------------")
//
//   for {
//     fmt.Print("-> ")
//     text, _ := reader.ReadString(' ')
//     fmt.Println(text)
//     // convert CRLF to LF
//     text = strings.Replace(text, "\n", "", -1)
//
//     if strings.Compare("hi", text) == 0 {
//       fmt.Println("hello, Yourself")
//     }
//
//   }
//
// }

// [_Command-line arguments_](http://en.wikipedia.org/wiki/Command-line_interface#Arguments)
// are a common way to parameterize execution of programs.
// For example, `go run hello.go` uses `run` and
// `hello.go` arguments to the `go` program.

package main

import "os"
import "fmt"

func main() {

    // `os.Args` provides access to raw command-line
    // arguments. Note that the first value in this slice
    // is the path to the program, and `os.Args[1:]`
    // holds the arguments to the program.
    argsWithProg := os.Args
    argsWithoutProg := os.Args[1:]

    // You can get individual args with normal indexing.
    arg := os.Args[3]

    fmt.Println(argsWithProg)
    fmt.Println(argsWithoutProg)
    fmt.Println(arg)
}
