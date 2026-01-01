package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) WriteDecrypt(encKey []byte, key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransfromFunc(key)
	fmt.Println("pathKey:pathName => ", pathKey.Pathname)
	fmt.Println("pathKey:Original => ", pathKey.Original)

	pathNameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())

	f, err := os.Create(pathAndFilenameWithRoot)
	if err != nil {
		return 0, err
	}
	n, err := copyDecrypt(encKey, r, f)
	if err != nil {
		return 0, err
	}

	//	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)
	return int64(n), nil
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransfromFunc(key)
	fmt.Println("pathKey:pathName => ", pathKey.Pathname)
	fmt.Println("pathKey:Original => ", pathKey.Original)

	pathNameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())

	f, err := os.Create(pathAndFilenameWithRoot)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	//	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)
	return n, nil
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)

}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s", s.Root+"/"+pathKey.Filename())
	file, err := os.Open(pathAndFilenameWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Filename())
	_, err := os.Stat(pathAndFilenameWithRoot)
	return !errors.Is(err, os.ErrNotExist)

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

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
