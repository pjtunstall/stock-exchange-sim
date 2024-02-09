package main

import (
	"fmt"
	"time"
)

func setTimer(f float64, c chan<- struct{}) {
	timer := time.After(time.Duration(f * float64(time.Second)))
	go wait(timer, c)
}

func wait(timer <-chan time.Time, c chan<- struct{}) {
	<-timer
	fmt.Println("Timer expired")
	c <- struct{}{}
}
