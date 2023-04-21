package main

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "os"
   "strings"
   "sync"
   "time"
)

func main() {
    start := time.Now()
    var wg sync.WaitGroup
    for i := 1; i <= 100; i++ {
        url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", i)
        wg.Add(1)
        go func(url string) {
        defer wg.Done()
        getEmails(url)
        }(url)
    }
    wg.Wait()
    end := time.Now()
    fmt.Printf("Elapsed time: %s", end.Sub(start))
}

func getEmails(url string) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    emails := extractEmails(string(body))
    saveEmails(emails)
}

func extractEmails(body string) []string {
    var emails []string
    lines := strings.Split(body, "\n")
    for _, line := range lines {
        if strings.Contains(line, "email") {
        email := strings.Trim(strings.Split(line, ":")[1], "\", ")
        emails = append(emails, email)
        }
    }
    return emails
}

func saveEmails(emails []string) {
    f, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    for _, email := range emails {
        _, err := f.WriteString(email + "\n")
        if err != nil {
        fmt.Println(err)
        return
        }
    }
}