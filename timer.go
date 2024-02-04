package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func setTimer() {
	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("error parsing wait time: %v", err)
	}
	timer := time.After(time.Duration(n) * time.Second)
	go wait(timer)
}

func wait(timer <-chan time.Time) {
	<-timer
	fmt.Println("Timer expired")
	os.Exit(0)
}
