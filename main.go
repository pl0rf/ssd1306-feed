package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fake := flag.Bool("fake", false, "use stdout as fake device")
	flag.Parse()

	go func() {
		s := &Stats{}
		for {
			s.LoadStats()
			s.DisplayStats()
			time.Sleep(5 * time.Second)
		}
	}()

	exit := make(chan struct{})
	output := make(chan string, 100)

	if *fake {
		go fakeDeviceLoop(output, exit)
	} else {
		go deviceLoop(output, exit)
	}

	c := NewClient(output)
	go c.Run()

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("interrupt")
	close(exit)
	c.Close()
}
