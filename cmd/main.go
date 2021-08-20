package main

import (
	"github.com/poncheska/discord-timetable/internal/app"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	go func() {
		r := http.NewServeMux()
		r.HandleFunc("/status", statusHandler)

		err := http.ListenAndServe(":8081", r)
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	app.Start()
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
