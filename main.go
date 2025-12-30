package main

import (
	"log"

	"github.com/imim77/gofilestorage/p2p"
)

func main() {
	trOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		ShakeHands: p2p.NOPHandShakefunc,
		Decoder:    p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(trOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}
