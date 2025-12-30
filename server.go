package main

import (
	"fmt"
	"log"

	"github.com/imim77/gofilestorage/p2p"
)

type FileServerOptions struct {
	listenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransfromFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOptions
	store  *Store
	quitch chan struct{}
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
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
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

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.loop()

	return nil

}
