package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "ggnetwork"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashString) / blockSize

	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Original: hashString,
	}

}

type PathTransfromFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Original string
}

func (p PathKey) Filename() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Original)
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Original: key,
	}
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folders
	// files of the system
	Root              string
	PathTransfromFunc PathTransfromFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransfromFunc == nil {
		opts.PathTransfromFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransfromFunc(key)
	fmt.Println("pathKey:pathName => ", pathKey.Pathname)
	fmt.Println("pathKey:Original => ", pathKey.Original)

	pathNameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil
	}

	pathAndFilename := pathKey.Filename()
	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())

	f, err := os.Create(pathAndFilenameWithRoot)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)
	return nil
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err

}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransfromFunc(key)
	//f, err := os.Open("d9e06/924cb/e4f7c/5f592/69e62/67f97/1d027/74564/d9e06924cbe4f7c5f59269e6267f971d02774564")
	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())
	f, err := os.Open(pathAndFilenameWithRoot)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())
	_, err := os.Stat(pathAndFilenameWithRoot)
	if err == fs.ErrNotExist {
		return false
	}
	return true

}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Filename())
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Original)
	}()
	seperatedFilePath := strings.Split(pathAndFilenameWithRoot, "/")

	return os.RemoveAll(seperatedFilePath[0])
}
