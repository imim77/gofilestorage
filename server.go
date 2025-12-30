package main

import (
	"fmt"
	"log"
	"net"
	"sync"

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
	peers    map[net.Addr]p2p.Peer
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
		peers:             make(map[net.Addr]p2p.Peer),
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr()] = p
	log.Printf("connected with remote %s", p.RemoteAddr())
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()

	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:

			return
		}

	}
}

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
