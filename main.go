package main

import (
	"log"
	"time"

	"github.com/imim77/gofilestorage/p2p"
)

func TestFunc(p p2p.Peer) error {
	p.Close()
	return nil
}

func main() {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		ShakeHands: p2p.NOPHandShakefunc,
		Decoder:    p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileserverOpts := FileServerOptions{

		StorageRoot:       ":3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	s := NewServer(fileserverOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
