package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "Foo not bar"
	src := bytes.NewReader([]byte(payload))
	dest := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dest)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dest.String())

	fmt.Println(len(payload))

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dest, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {
		t.Fail()
	}

	if out.String() != payload {
		t.Errorf("decryption failed: %s", err)
	}

	fmt.Println(out.String())
}
