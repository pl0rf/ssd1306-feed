package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ws        *websocket.Conn
	output    chan string
	connected bool
	err       error
	quit      chan struct{}
	done      chan struct{}

	// TODO: add some stats on the ws client?
	// reconnects int
	// created    time.Time
	// connected  time.Time
	// recvRate   int
}

func NewClient(output chan string) *Client {
	c := &Client{
		output: output,
		quit:   make(chan struct{}),
	}
	return c
}

func (c *Client) Close() {
	log.Println("ws client close")
	close(c.quit)

	if c.connected { // XXX: race condition here, but who cares
		err := c.ws.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.ws.Close()
		if err != nil {
			log.Println("write close:", err)
			return
		}
	}

	select {
	case <-c.done:
	case <-time.After(time.Second):
	}

	return
}

func (c *Client) Run() {
	for {
		log.Println("ws client attempt to connect")
		c.output <- "DIAL"
		c.ws, _, c.err = websocket.DefaultDialer.Dial(COINBASE_URL, nil)
		if c.err != nil {
			log.Println(c.err)
			c.output <- "ERROR"
			select {
			case <-c.quit:
				log.Println("ws client quitting")
				return
			default: // reconnect timeout
				time.Sleep(3 * time.Second)
			}
			continue
		}
		log.Println("ws client connected")
		c.done = make(chan struct{})
		c.connected = true
		c.readLoop()
		c.connected = false
		c.output <- "ERROR"
		select {
		case <-c.quit:
			return
		default:
		}
	}
}

func (c *Client) readLoop() {
	defer func() {
		close(c.done)
	}()

	if err := c.ws.WriteJSON(subscribe); err != nil {
		log.Println("write msg error: ", err)
		return
	}

	var msg Message

	for {
		msg = Message{}
		if err := c.ws.ReadJSON(&msg); err != nil {
			log.Println("read msg error: ", err)
			return
		}

		if msg.Type == "ticker" {
			select {
			case c.output <- formatStr(msg.Price):
			default:
				log.Println("Channel full. Discarding msg.")
				<-c.output // XXX: this is probably a bad idea
				c.output <- msg.Price
			}

		} else {
			log.Println(msg.Type)
		}

	}
}
