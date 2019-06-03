package main

import (
    "bufio"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "fmt"
    "regexp"
    "encoding/csv"
    "strconv"
    //"errors"
)

func notPresentIn(keys []string, k string) bool{
  for _, e := range keys{
    if e == k {return false}
  }
  return true
}
func sliceUnion(keys, ks []string) []string{
  for _, e := range ks{
    if notPresentIn(keys, e){
      keys = append(keys, e)
    }
  }
  return keys
}
func check(e error) {
    if e != nil {
        log.Fatal(e)
    }
}
func writeInFile(s string, fileName string) {
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
func main(){
  rs := links()
  fmt.Println(len(rs))
  keys := make([]string, 0)
  for i, s := range rs{
    key := mainHere(s, i)
    keys = sliceUnion(keys, key)
    fmt.Println("Done for "+s)
  }
  fmt.Println("___*****____")
  fmt.Println("Count : ", len(keys))
  printSlice(keys)
}
func mainHere(urll string, i int)[]string{
  response, err := http.Get(urll)
  check(err)
  defer response.Body.Close()
  dataInBytes, err := ioutil.ReadAll(response.Body)
  check(err)
  pageContent := string(dataInBytes)
  mapp := regexData(returnFeeInfo(pageContent, urll))
  keys := make([]string, 0)
  for k, _ := range mapp{
    keys = append(keys, k)
  }
  //writeToCSV(mapp, i)
  return keys
}
func printSlice(m []string){
  for _, v := range m{
    fmt.Println(v)
  }
}
func printMap(m map[string]string){
  for k, v := range m{
    fmt.Println(k+" -> "+v)
  }
}
func returnFeeInfo(pageContent, urll string) string{
  var rdx, rdx2 int
  rdx = strings.Index(pageContent, "class=\"tab-inner-content\" id=\"feature-2-tab\"")
  if rdx <0 {
    fmt.Println("Start", urll)
    return ""
    } else{
    rdx2 = strings.Index(pageContent[rdx:], "</section>")
    if rdx2 <0 {
      fmt.Println("End", urll)
      return ""
    }
  }
  return pageContent[rdx:rdx2+rdx]
}

func regexData(pageC string) map[string]string {
  mapp := make(map[string]string)
  re := regexp.MustCompile("<li>([^<:]+)(?:[ \t]*):(?:[ \t]*)(?:<em(?:[ \t]*)class=\"WebRupee\">(?:[ \t]*)Rs.(?:[ \t]*)</em>)?(.*)</li>")
  comments := re.FindAllStringSubmatchIndex(pageC, -1)
  //fmt.Println(comments)
  if comments == nil {
        fmt.Println("No matches for table.")
        writeInFile(pageC, "test.html")
    } else {
        for _, comment := range comments {
            mapp[pageC[comment[2]:comment[3]]] = pageC[comment[4]:comment[5]]
        }
    }
    return mapp
}
func writeToCSV(m map[string]string, i int){
  file, err := os.Create("SBI/SBIresult"+strconv.Itoa(i)+".csv")
  check(err)
  defer file.Close()
  writer := csv.NewWriter(file)
  defer writer.Flush()
  er := writer.Write([]string{"Fee Type", "Value"})
  check(er)
  for k, v := range m{
    err := writer.Write([]string{k, v})
    check(err)
  }
  writer.Flush()
  fmt.Println("Wrote to CSV")
}

func links() []string{
  rs := make([]string, 0)
  file, err := ioutil.ReadFile("sbilinks.html")
  check(err)
  //defer file.Close()
  pageC := string(file)
  re:=regexp.MustCompile("a(?:[\t ]*)href=\"([^\"]+)\"(?:[\t ]*)class=\"(?:[\t ]*)learn-more-link(?:[\t ]*)\"")
  cc := re.FindAllStringSubmatchIndex(pageC, -1)
  if cc == nil{
    fmt.Println("Nill")
  } else{
    for _, c := range cc{
      rs = append(rs, pageC[c[2]:c[3]])
    }
  }
  return rs
}
