package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/redtoad/go-pibrella/pibrella"
)

func main() {

	// clean exit on ctrl + c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			//yellow.Low()
			//green.Low()
			//red.Low()
			_ = pibrella.Close()
			os.Exit(0)
		}
	}()

	fmt.Print("Hello world!\n")

	err := pibrella.Open()
	if err != nil {
		panic(fmt.Errorf("could not open GPIO: %v", err))
	}
	defer pibrella.Close()

	pibrella.Yellow.On()
	time.Sleep(2 * time.Second)

	for {
		pibrella.Green.Toggle()
		pibrella.Yellow.Toggle()
		pibrella.Red.Toggle()

		time.Sleep(time.Second)
	}

}
