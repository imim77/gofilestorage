package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)
	//expectedPathName := "d9e06/924cb/e4f7c/5f592/69e62/67f97/1d027/74564/d9e06924cbe4f7c5f59269e6267f971d02774564"

}

func TestOnlyDelete(t *testing.T) {
	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	s := NewStore(opt)
	key := "mombestpicture"
	if err := s.Delete(key); err != nil {
		t.Error(t)
	}
}

func TestOnlyHas(t *testing.T) {

	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	s := NewStore(opt)

	key := "mombestpicture"

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

}

func TestStoreDelete(t *testing.T) {
	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	s := NewStore(opt)

	key := "mombestpicture"

	data := []byte("some jpeg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	s := NewStore(opt)

	key := "mombestpicture"

	data := []byte("some jpeg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s, have %s", data, b)
	}
}
