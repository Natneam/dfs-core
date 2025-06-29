package cipher

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	data := []byte("Hello world")
	src := bytes.NewReader(data)
	dst := new(bytes.Buffer)
	key := NewEncryptionKey()
	if _, err := CopyEncrypt(key, src, dst); err != nil {
		t.Fail()
	}

	out := new(bytes.Buffer)

	if _, err := CopyDecrypt(key, dst, out); err != nil {
		t.Fail()
	}

	assert.Equal(t, data, out.Bytes())
}
