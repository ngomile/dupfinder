package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"os/path"
)

var (
	dir string
)

func crc32Hash(filepath string) (uint32, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	h := crc32.NewIEEE()
	h.Write(b)
	return h.Sum32(), nil
}

func main() {
	flag.StringVar(&dir, "dir", "", "Specifies the directory to search")
	flag.Parse()

	if dir == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}	

	matches := make(map[uint32][]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	
	var (
		fpath string
		crc32h uint32
	)

	for file := range files {
		fpath = path.Join(dir, file.Name())
		crc32h, err = crc32Hash(fpath)
		
		if err != nil {
			log.Print("Failed to retrieve hash of:", fpath)
			continue
		}

		if match, ok := matches[crc32h]; ok {
			matches[crc32h] = append(matches[crc32h], fpath)
		} else {
			matches[crc32h] = []string{fpath}
		}
	}
}