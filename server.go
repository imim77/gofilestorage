package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/imim77/gofilestorage/p2p"
)

type FileServerOptions struct {
	listenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransfromFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOptions

	peerLock sync.RWMutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewServer(opts FileServerOptions) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransfromFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		store:             NewStore(storeOpts),
		FileServerOptions: opts,
		quitch:            make(chan struct{}),
		peers:             make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddr())
	return nil
}

func (s *FileServer) Broadcast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)

}

type Message struct {
	Payload any
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	//buf := new(bytes.Buffer)
	//tee := io.TeeReader(r, buf)
	buf := new(bytes.Buffer)
	msg := Message{
		Payload: []byte("storeage key"),
	}
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	time.Sleep(time.Second * 3)
	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil

	//if err := s.store.Write(key, tee); err != nil {
	//	return err
	//}

	//p := &DataMessage{
	//	Key:  key,
	//	Data: buf.Bytes(),
	//}
	//fmt.Println(p.Data)
	//return s.Broadcast(&Message{
	//	From:    "todo",
	//	Payload: p,
	//})
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()

	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():

			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
			}
			fmt.Printf("recieved: %s\n", string(msg.Payload.([]byte)))

			peer, ok := s.peers[rpc.From]
			if !ok {
				panic("peer not founc in peers map")
			}
			b := make([]byte, 1000)

			if _, err := peer.Read(b); err != nil {
				panic(err)
			}

			fmt.Printf("%s\n", string(b))
			peer.(*p2p.TCPPeer).Wg.Done()
			//if err := s.handleMessage(&m); err != nil {
			//	log.Panicln(err)
			//}

		case <-s.quitch:

			return
		}

	}
}

//func (s *FileServer) handleMessage(msg *Message) error {
//	switch v := msg.Payload.(type) {
//	case *DataMessage:
//		fmt.Printf("recieved data %+v\n", v)
//	}
//	return nil
//}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func() {
			fmt.Println("attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)

			}
		}()

	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil

}
