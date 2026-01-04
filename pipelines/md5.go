package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// Example usage of the MD5 pipeline
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println("err while getting md5sum", err)
		return
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	// Sort paths for consistent output
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

// MD5All computes the MD5 checksum for all files in the specified directory.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil // skip non-regular files
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
