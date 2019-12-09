package main

import (
	"log"
	"time"

	"github.com/jjcinaz/wsClient1/pixelblaze"
)

func main() {
	var (
		err error
		pb  *pixelblaze.Pixelblaze
	)
	log.SetFlags(0)
	pb, err = pixelblaze.Connect("192.168.91.140:81")
	if err != nil {
		log.Fatal(err)
	}
	pb.Write(`{"getConfig":true,"listPrograms":true,"sendUpdates":false}`)
	time.Sleep(time.Second)
	pb.Close(time.Second)
}
