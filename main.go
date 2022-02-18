package main

import (
	"log"
	"logsCollector/internal/db"
	"logsCollector/internal/insertlogs"

	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func init() {
	log.Println("init", "started")
	defer log.Println("init", "finished")

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	c := db.Connect()
	log.Println("init", "connect ", c)

}

func main() {

	//http.HandleFunc("/inslog", insertlogs.InsLogs)
	http.HandleFunc("/inslog", insertlogs.InsLogs)

	http.ListenAndServe(":3021", nil)
}
