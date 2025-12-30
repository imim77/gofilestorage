package main

import (
	"fmt"
	"log"

	"github.com/imim77/gofilestorage/p2p"
)

func TestFunc(p p2p.Peer) error {
	p.Close()
	return nil
}

func main() {
	trOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		ShakeHands: p2p.NOPHandShakefunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer:     TestFunc,
	}
	tr := p2p.NewTCPTransport(trOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}
