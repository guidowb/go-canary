package main

import (
	"os"
	"os/signal"
	"time"
	"syscall"
	"net/http"
	"log"
	"strconv"
	"strings"
)

func main() {

    logMessage("Starting")

    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)

    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	go func() {
        <-sigs
        logMessage("Initiating graceful shutdown")
        countdown := 19
        for countdown > 0 {
            logMessage("Shutting down in " + strconv.Itoa(countdown) + " seconds")
            time.Sleep(1000 * time.Millisecond)
            countdown--
        }
	    done <- true
    }()

    <-done
    logMessage("Shutting down now")
}

func logMessage(message string) {
	http.Post("http://shutdown-logger.cfapps.io", "text/xml", strings.NewReader(message))
}