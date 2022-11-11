package encrypt_test

import (
	"Lockr/pkg/encrypt"
	"io"
	"os"
	"testing"
)

func TestStreamEncrypt(t *testing.T) {
	file, err := os.Open("testfile.txt")
	if err != nil {
		t.Fatal("failed to find test target file")
	}

	r, err := encrypt.NewEncryptReader(file, []byte("hello"))
	if err != nil {
		t.Fatal("failed to create new encrypt reader")
	}
	endpoint, err := os.Create("testoutput.txt")
	if err != nil {
		t.Fatal("failed to open testoutput")
	}
	io.Copy(endpoint, r)
}

func TestStreamDecrypt(t *testing.T) {
	file, err := os.Open("testoutput.enc")
	if err != nil {
		t.Fatal("failed to find test target file")
	}

	r, err := encrypt.NewDecryptReader(file, []byte("hello"))
	if err != nil {
		t.Fatal("failed to create new decrypt reader")
	}

	endpoint, err := os.Create("decryptedoutput.txt")
	if err != nil {
		t.Fatal("failed to open decryptedoutput")
	}
	_, err = io.Copy(endpoint, r)
	if err != nil {
		t.Fatal(err)
	}

}
