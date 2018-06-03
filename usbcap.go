package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/gousb"
)

func main() {
	// following https://godoc.org/github.com/google/gousb#Context.OpenDeviceWithVIDPID
	ctx := gousb.NewContext()
	defer ctx.Close()

	ctx.Debug(4)

	dev, err := ctx.OpenDeviceWithVIDPID(0x048d, 0x1234)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	}
	defer done()

	ep, err := intf.OutEndpoint(7)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(7): %v", intf, err)
	}

	bytesWritten := 0
	bytesAtATime := 100 * 1024 * 1024 // 100 MB
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, bytesAtATime)
	for {
		rand.Read(data)
		numBytes, err := ep.Write(data)
		if err != nil {
			break
		}
		bytesWritten += numBytes
		fmt.Printf("%d bytes written\n", bytesWritten)
		if numBytes != bytesAtATime {
			break
		}
	}
	fmt.Printf("Done! %d bytes written!\n", bytesWritten)
}
