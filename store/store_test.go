package store

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "test49"
	pathKey := HashPathTransformFunc(key)
	expectedFileName := "241db62d4cb712d490d6c6fdd81c4682"
	expectedFolderName := "241db6/2d4cb7/12d490/d6c6fd/d81c46"

	assert.Equal(t, pathKey.FileName, expectedFileName)
	assert.Equal(t, pathKey.PathName, expectedFolderName)
}

func TestWrite(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := range 50 {
		key := fmt.Sprintf("test%d", i)

		if err := s.Write(key, bytes.NewReader([]byte("Hello World"))); err != nil {
			t.Error(err)
		}

		if !s.Has(key) {
			t.Errorf("expected to have key %s", key)
		}
	}
}

func TestRead(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := range 50 {
		key := fmt.Sprintf("test%d", i)

		if err := s.Write(key, bytes.NewReader([]byte("Hello World"))); err != nil {
			t.Error(err)
		}

		data, err := s.Read(key)
		if err != nil {
			t.Error()
		}

		assert.Equal(t, data, []byte("Hello World"))
	}
}

func TestDelete(t *testing.T) {
	s := newStore()
	defer tearDown(t, s)

	for i := range 50 {
		key := fmt.Sprintf("test%d", i)

		if err := s.Write(key, bytes.NewReader([]byte("Hello World"))); err != nil {
			t.Error()
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if s.Has(key) {
			t.Errorf("expected to not have key %s", key)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: HashPathTransformFunc,
	}
	return NewStore(opts)
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
