package main

import "fmt"

type Reminder struct {
}

func (r Reminder) Run() {
	fmt.Printf("Every 5 sec remind\n")
}
