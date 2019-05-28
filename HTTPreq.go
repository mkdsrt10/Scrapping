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
    //"errors"
)
type TupleStr struct {
  key string
  value string
}
func notPresentIn(keys []string, k string) bool{
  for _, e := range keys{
    if e == k {return false}
  }
  return true
}
func check(e error) {
    if e != nil {
        panic(e)
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
func FeeIndex(pageContent string) int{
  searchName := []string{"Fees and Charges", "Fees And Charges", "Fees & Charges", "Fees and charges", "Fee &amp; Charges", "Fees &amp; Charges", "fees and charges", "Charges for" }
  for _, el := range searchName {
    i := strings.Index(pageContent, el)
    // fmt.Println(i)
    if i >= 0{
      return i
    }
  }
  // writeInFile(pageContent, "tabletest.html")
  // panic(errors.New("Can't find Fees and Charges"))
  return -10
}
func returnFeeTable(pageContent string) []TupleStr{
  table := make([]TupleStr, 0)
  //writeInFile(pageContent, "writeTable.html")
  //Create a regular expression to find comments
  re:= regexp.MustCompile("<tr>(?:\r\n?)<td(?:[\t\r\n ]*colspan=\"[0-9]\"[\t\r\n ]*)?>(?:[\t\n\r ]*)?(?:<p>)?(?:<strong>)?([^</]+)(?:</strong>)?(?:</p>)?(?:[\n\t\r ]*)</td>(?:\r\n?)<td(?:[\t\r\n ]*colspan=\"[0-9]\"[\t\r\n ]*)?>(?:[\t\n\r ]*)?(?:<p>)?(?:<strong>)?([^</]+)(?:</strong>)?(?:</p>)?(?:[\n\t\r ]*)</td>(?:\r\n?)</tr>")
  comments := re.FindAllStringSubmatchIndex(pageContent, -1)
  //fmt.Println(comments)
  if comments == nil {
        fmt.Println("No matches for table.")
        writeInFile(pageContent, "test.html")
    } else {
        for _, comment := range comments {
            typeFee, amountS := pageContent[comment[2]:comment[3]], pageContent[comment[4]:comment[5]]
            table = append(table, TupleStr{typeFee, amountS})
        }
    }
    return table
}
func returnNameCard(pageContent string) TupleStr{
  //writeInFile(pageContent, "writeTable.html")
  //Create a regular expression to find comments
  re:= regexp.MustCompile("<meta property=\"og:title\" content=\"([^-,:\"]*)")
  comments := re.FindStringSubmatch(pageContent)
  //fmt.Println(comments)
  if comments == nil {
        fmt.Println("No matches for name.")
    } else {
      return TupleStr{"name of card", comments[1]}
    }
    return TupleStr{"name of card", ""}
}
func returnHiddenF(pageContent string) []TupleStr{
  var table []TupleStr
  re:= regexp.MustCompile("<input type=\"hidden\" name=\"featureTypeName\" class=\"(?:.*)\" data-(.*)=\"(.*)\" />")
  comments := re.FindAllStringSubmatchIndex(pageContent, -1)
  //fmt.Println(comments)
  if comments == nil {
        fmt.Println("No matches for Hidden data.")
    } else {
      for _, comment := range comments {
          typeFee, amountS := pageContent[comment[2]:comment[3]], pageContent[comment[4]:comment[5]]
          table = append(table, TupleStr{typeFee, amountS})
      }
    }
    return table
}
func returnCardData(urll string) []TupleStr{
    // Make HTTP GET request
    // fmt.Println("Start scrapping from ", urll)
    var hiddenContent string
    var table []TupleStr
    var titleStartIndex, titleEndIndex int
    var pageContent string
    response, err := http.Get(urll)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    // fmt.Println("Reached the page")
    // Get the response body as a string
    dataInBytes, err := ioutil.ReadAll(response.Body)
    check(err)
    pageContent = string(dataInBytes)
    // writeInFile(pageContent, "writeFile.html")
    titleStartIndex = strings.Index(pageContent, "<head")
    titleEndIndex = strings.Index(pageContent, "</head>")
    headContent := pageContent[titleStartIndex:titleEndIndex]
    table = append(table, returnNameCard(headContent))
    //fmt.Println(name)
    t := strings.Index(pageContent, "<body")
    titleStartIndex = FeeIndex(pageContent[t:])
    if titleStartIndex != -10 {
      tableContent := pageContent[titleStartIndex+t:]
      titleEndIndex = strings.Index(tableContent, "</table>")
      if titleEndIndex >= 0{
        tableContent = tableContent[:titleEndIndex]
      }
      for _,elm := range returnFeeTable(tableContent){
        table = append(table, elm)
      }
    }
    // table = append(table, returnFeeTable(tableContent))
    //fmt.Println(table)
    titleStartIndex = strings.Index(pageContent, "offers-section  offer")
    if titleStartIndex >= 0{
      hiddenContent = pageContent[titleStartIndex:]
      titleEndIndex = strings.Index(hiddenContent, "<div class=\"js-offer-compare-table\"></div>")
      if titleEndIndex >= 0{
        hiddenContent = hiddenContent[:titleEndIndex]
      }
    }
    //writeInFile(hiddenContent, "HiddenC.html")
    for _,elm := range returnHiddenF(hiddenContent){
      table = append(table, elm)
    }
    //table = append(table, returnHiddenF(hiddenContent))
    // fmt.Println("Scrapped data from ", urll)
    return table
}
func urlList(urll string) []string{
  // <a  href="(\S*) "  class="js-title js-format-string value-title ">
  // Make HTTP GET request
  fmt.Println("Start scrapping URL list from ", urll)
  var table []string
  // var titleStartIndex, titleEndIndex int
  var pageContent string
  response, err := http.Get(urll)
  if err != nil {
      log.Fatal(err)
  }
  defer response.Body.Close()
  fmt.Println("Reached the page")
  // Get the response body as a string
  dataInBytes, err := ioutil.ReadAll(response.Body)
  check(err)
  pageContent = string(dataInBytes)

  re:= regexp.MustCompile("<a  href=\"([^\t\n\f\r ]*) \"  class=\"js-title js-format-string value-title \">")
  comments := re.FindAllStringSubmatchIndex(pageContent, -1)
  //fmt.Println(comments)
  if comments == nil {
        fmt.Println("No matches.")
    } else {
        for _, comment := range comments {
            table = append(table, pageContent[comment[2]:comment[3]])
        }
    }
    return table
}
func main(){
  if os.Args[1] == "multi"{
    // var keys []string
    siteList := links()
    // var keyVal [][]TupleStr

    file, err := os.Create("result.csv")
    check(err)
    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()
    er := writer.Write([]string{"Name of credit card","Bank","Network","Annual Fee", "Annual Fee Waiver","Joining Fee","Interest Rate","Cash Withdrawal Fee","Cash Advance Interest","Foreign currency transaction fee","Fuel Surcharge"})
    check(er)

    for _, r := range siteList{
      fmt.Println("____Start all Scrapping____\n")
      urllist := urlList("https://www.bankbazaar.com"+r)
      fmt.Println("Site : ", r)
      fmt.Println("No. of cards :", len(urllist))
      fmt.Println("\n________*******________\n")
      for _, elm := range urllist{
        err := writer.Write(toOurFormat(toMap(returnCardData("https://www.bankbazaar.com"+elm))))
        check(err)
        // keyVal = append(keyVal, returnCardData("https://www.bankbazaar.com"+elm))
      }
    }
    fmt.Println("____Done all Scrapping____\n")
    defer writer.Flush()
    // writeToCSV(keyVal)
  } else {
    a := returnCardData(os.Args[1])
    fmt.Println(toMap(a))
      // for i, el := range a{
      //   fmt.Println(string(i)+" : "+el.key+" -> "+el.value)
      // }
  }
}


func toOurKey(k string) string{
  switch false{
  case k != "name of card":
    return "Name of Card"
  case notPresentIn([]string{"annual-fee", "First Year Annual Fee",	"Annual fee (from 2nd of card membership)",	"Annual fee from year-2 onwards (primary cardholder)", "Annual fee", "First year annual fee",	"Annual fees"}, k):
    return "Annual Fee"
  case notPresentIn([]string{"Joining Fee (1st year)", "Joining fee",	"Joining Fee (Primary Cardholder)",	"Joining fees"}, k):
    return "Joining Fee"
  case notPresentIn([]string{"Minimum Spends for Annual Fee Reversal", "Minimum spend for waiver of annual fee"}, k):
    return "Annual Fee Waiver"
  case notPresentIn([]string{"Overdue interest in extended credit", "Charges on purchases",	"Interest rate (cash and retail purchases)",	"Finance charges",	"Finance charges - cash and retail transactions", "Finance charges (cash and retail purchases)"}, k):
    return "Interest Rate"
  case notPresentIn([]string{"Cash advance charge", "Cash withdrawal charges",	"Cash withdrawal or cash advance fees",	"Fees on cash transaction", "Finance charge for cash advance",	"Finance charges - on cash advances",	"Interest on cash advances"}, k):
    return "Cash Advance Interest"
  case notPresentIn([]string{"Cash withdrawal fee", "Fee for cash withdrawal"}, k):
    return "Cash Advance Fee"
  case notPresentIn([]string{"Payments in foreign currency",	"Foreign currency transaction fee",	"Foreign currency transaction charge",	"Foreign Currency Transaction Fee",	"Fee for foreign currency transaction"}, k):
    return "Foreign Currency Mark-up"
  case k != "networks":
    return "Network"
  case k != "bank":
    return "Bank"
  case k != "Surcharge on fuel-related payments":
    return "Fuel Surcharge"
  default :
    return k
  }
}
func toMap(tup []TupleStr) map[string]string{
  mapp := make(map[string]string)
  for _, el := range tup{
    mapp[toOurKey(el.key)] = el.value
  }
  return mapp
}
func toOurFormat(keyVal map[string]string)[]string{
  var result []string
  result = append(result, keyVal["Name of Card"])
  result = append(result, keyVal["Bank"])
  result = append(result, keyVal["Network"])
  result = append(result, keyVal["Annual Fee"])
  result = append(result, keyVal["Annual Fee Waiver"])
  result = append(result, keyVal["Joining Fee"])
  result = append(result, keyVal["Interest Rate"])
  result = append(result, keyVal["Cash Advance Fee"])
  result = append(result, keyVal["Cash Advance Interest"])
  result = append(result, keyVal["Foreign Currency Mark-up"])
  result = append(result, keyVal["Fuel Surcharge"])
  for k, e := range keyVal{
    if notPresentIn([]string{"Name of Card","Bank","Network", "Annual Fee", "Joining Fee", "Interest Rate", "Cash Advance Fee", "Cash Advance Interest", "Foreign Currency Mark-up", "Fuel Surcharge",    "joining-perks", "popularity", "card-fee-type"}, k){
      result = append(result, k+" : "+e)
    }
  }
  return result
}
// func writeToCSV(keyVal [][]TupleStr){
//   file, err := os.Create("result.csv")
//   check(err)
//   defer file.Close()
//   writer := csv.NewWriter(file)
//   defer writer.Flush()
//   er := writer.Write([]string{"Name of credit card","Bank","Network","Annual Fee","Joining Fee","Interest Rate","Cash Withdrawal Fee","Cash Advance Interest","Foreign currency transaction fee","Fuel Surcharge"})
//   check(er)
//   for _, value := range keyVal {
//       err := writer.Write(toOurFormat(toMap(value)))
//       check(err)
//   }
// }


func links() []string {
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
