// find_html_comments_with_regex.go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
)

func main() {
    // Make HTTP request

    // Read response data in to memory
    dat, err := ioutil.ReadFile("links.txt")
    check(err)
    cont := string(dat)
    resl := make([]string, 0)
    // Create a regular expression to find comments
    re := regexp.MustCompile("<li><a title=\"(?:[^\"]*)\" href=\"([^\"]*)\"")
    comments := re.FindAllStringSubmatchIndex(cont, -1)
    if comments == nil {
        fmt.Println("No matches.")
    } else {
        for _, comment := range comments {
          resl = append(resl, cont[comment[2]:comment[3]])
        }
    }
    return resl
}
