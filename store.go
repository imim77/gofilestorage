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
	Root string
	// ID of the owner of the storage, which will be used to store all
	//files at that location so we can sync all the files if needed.
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

func (s *Store) Write(id, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

func (s *Store) WriteDecrypt(encKey []byte, id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
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

func (s *Store) openFileForWriting(id, key string) (*os.File, error) {
	pathKey := s.PathTransfromFunc(key)
	fmt.Println("pathKey:pathName => ", pathKey.Pathname)
	fmt.Println("pathKey:Original => ", pathKey.Original)

	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}

	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Filename())

	return os.Create(pathAndFilenameWithRoot)
}

func (s *Store) writeStream(id, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)

}

func (s *Store) Read(id, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)

}

func (s *Store) readStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Filename())

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

func (s *Store) Has(id, key string) bool {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Filename())
	_, err := os.Stat(pathAndFilenameWithRoot)
	return !errors.Is(err, os.ErrNotExist)

}

func (s *Store) Delete(id, key string) error {
	pathKey := s.PathTransfromFunc(key)
	pathAndFilenameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Filename())
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Original)
	}()
	seperatedFilePath := strings.Split(pathAndFilenameWithRoot, "/")

	return os.RemoveAll(seperatedFilePath[0])
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
