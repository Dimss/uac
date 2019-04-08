package main

import (
	"fmt"
	"time"
)

func chTest(myChan chan string) {
	fmt.Println("Waiting for channel input")
	fmt.Println(<-myChan)
	fmt.Println("Finished writes to channel, gonna exit routine ")
}

func startChan(myChan chan string) {
	time.Sleep(2 * time.Second)
	myChan <- "Hello world"
}

func main() {
	myChan := make(chan string)
	go chTest(myChan)
	go startChan(myChan)
	fmt.Println("continue execution here")
	time.Sleep(600 * time.Second)
	//messages := make(chan string)
	//go func() { messages <- "ping" }()
	//msg := <-messages
	//fmt.Println(msg)
}
