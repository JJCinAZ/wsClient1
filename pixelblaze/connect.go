package pixelblaze

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Pixelblaze struct {
	address   string
	writeChan chan []byte
	readDone  chan struct{}
	interrupt chan os.Signal
	c         *websocket.Conn
}

func (pb *Pixelblaze) Close(waittime time.Duration) {
	select {
	case pb.interrupt <- os.Interrupt:
		select {
		case <-pb.readDone:
		case <-time.After(waittime):
		}
	default:
	}
}

func (pb *Pixelblaze) Write(msg string) {
	pb.writeChan <- []byte(msg)
}

func Connect(address string) (*Pixelblaze, error) {
	var (
		err error
	)
	pb := new(Pixelblaze)
	pb.interrupt = make(chan os.Signal, 1)
	signal.Notify(pb.interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	log.Printf("connecting to %s", u.String())

	pb.c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	pb.writeChan = make(chan []byte)
	pb.readDone = make(chan struct{})
	// Start read loop
	go func() {
		defer close(pb.readDone)
		for {
			_, message, err := pb.c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	// Start message pump
	go func() {
		defer pb.c.Close()
		for {
			select {
			case <-pb.readDone:
				return
			case msg := <-pb.writeChan:
				err := pb.c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-pb.interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := pb.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-pb.readDone:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}()
	return pb, nil
}
