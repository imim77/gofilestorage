package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	return NewStore(opt)

}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(t)
	}
}

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

func TestStoreDelete(t *testing.T) {
	opt := StoreOpts{
		PathTransfromFunc: CASPathTransformFunc,
	}
	s := NewStore(opt)

	key := "mombestpicture"

	data := []byte("some jpeg bytes")
	if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpeg bytes")
		if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		_, r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, _ := io.ReadAll(r)
		fmt.Println(string(b))
		if string(b) != string(data) {
			t.Errorf("want %s, have %s", data, b)
		}
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}
	}

}
