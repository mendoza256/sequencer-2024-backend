package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
    "go-sql-driver/mysql"
)

type Sequence struct {
    name string
    notation [8][16]string
}

func (s *Sequence) save() {
    fmt.Printf("Save sequence \n")
}

db, err := sql.Open("mysql", "root@localhost:3306@/sequencer-2024")
if err != nil {
	panic(err)
}
// See "Important settings" section.
db.SetConnMaxLifetime(time.Minute * 3)
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(10)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func saveSequenceHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Path[len("/sequence/save/"):]

    var notation [8][16]string
    err := json.NewDecoder(r.Body).Decode(&notation)
    if err != nil {
        http.Error(w, "Error decoding request body", http.StatusBadRequest)
        return
    }

    seq := &Sequence{name: name, notation: notation}
    seq.save()
    fmt.Printf("Save log: %s", name)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("Test log")
}

func main() {
    http.HandleFunc("/test/", testHandler)
    http.HandleFunc("/sequence/save/", saveSequenceHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
