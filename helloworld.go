package main

import (
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"
)

func main() {
    t := time.Now()
    text := "Hello World Rico, it's demo time " + t.Format("20060102150405") + " :-)"

    rev := os.Getenv("K_REVISION") // K_REVISION=helloworld-s824d
    if i := strings.LastIndex(rev, "-"); i > 0 {
        rev = rev[i+1:]
    }

    msg := fmt.Sprintf("%s: %s\n", rev, text)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Got request\n")
        time.Sleep(500 * time.Millisecond)
        fmt.Fprint(w, msg)
    })

    fmt.Printf("Listening on port 8080 (rev: %s)\n", rev)
    http.ListenAndServe(":8080", nil)
}
