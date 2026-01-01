package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/imim77/gofilestorage/p2p"
)

func TestFunc(p p2p.Peer) error {
	p.Close()
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
		ShakeHands: p2p.NOPHandShakefunc,
		Decoder:    p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileserverOpts := FileServerOptions{

		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
		EncKey:            newEncryptionKey(),
	}
	s := NewServer(fileserverOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s

}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":8000", "")
	//s3 := makeServer(":5000", ":3000")
	s3 := makeServer(":5001", ":3000", ":8000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(time.Millisecond * 500)
	go func() {
		log.Fatal(s2.Start())
	}()

	time.Sleep(time.Second * 2)
	go func() { log.Fatal(s3.Start()) }()
	//key := fmt.Sprintf("picture.png")
	//data := bytes.NewReader([]byte("my big data file here!"))
	//s2.Store(key, data)
	time.Sleep(time.Second * 2)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here!"))
		s3.Store(key, data)

		if err := s3.store.Delete(s3.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}

}
