package main

import (
	"bytes"
	"fmt"
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
	}
	s := NewServer(fileserverOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s

}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 1)
	go s2.Start()
	time.Sleep(time.Second * 1)
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("myprivatedata_%s", i)
		data := bytes.NewReader([]byte("my big data file here"))
		s2.Store(key, data)
		time.Sleep(5 * time.Millisecond)
	}

	//r, err := s2.Get("myprivatedata")
	////r, err := s2.Get("aaaa")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//_, err = io.ReadAll(r)
	//if err != nil {
	//	log.Fatal(err)
	//}
	////fmt.Println(string(b))

	select {}

}
