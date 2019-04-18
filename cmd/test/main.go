package main

import (
	"fmt"
	"time"
)


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
