package store

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Root is a root directory which will be used to store files in.
const root = "../storedfiles"

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func (p PathKey) BaseFolder() (string, error) {
	folders := strings.Split(p.PathName, "/")
	if len(folders) > 0 {
		return folders[0], nil
	}

	return "", fmt.Errorf("invalid folder structure")
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	PathTransformFunc
	Root string
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = root
	}
	return &Store{StoreOpts: opts}
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) ([]byte, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf.Bytes(), err
}

func (s *Store) Delete(key string) error {
	pathName := s.PathTransformFunc(key)

	folder, err := pathName.BaseFolder()
	if err != nil {
		return err
	}
	return os.RemoveAll(fmt.Sprintf("%s/%s", s.Root, folder))
}

func (s *Store) Has(key string) bool {
	pathName := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathName.PathName)

	_, err := os.Stat(pathNameWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathName := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathName.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathName.FullPath())

	f, err := os.Create(fullPathWithRoot)

	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)

	if err != nil {
		return 0, err
	}

	log.Printf("Writing Data %d to %s", n, fullPathWithRoot)

	return n, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathName := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathName.FullPath())

	f, err := os.Open(fullPathWithRoot)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

func HashPathTransformFunc(key string) PathKey {
	hash := md5.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])
	paths := []string{}

	folderNameLength := len(hashString) / 5

	for i := 0; i < len(hashString); i += folderNameLength {
		if i+folderNameLength < len(hashString) {
			paths = append(paths, hashString[i:i+folderNameLength])
		}
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashString,
	}
}
