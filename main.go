package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func crc32Hash(filepath string) (uint32, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	h := crc32.NewIEEE()
	h.Write(b)
	return h.Sum32(), nil
}

func findDuplicates(dir string) (map[uint32][]string, error) {
	matches := make(map[uint32][]string)
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	var (
		fpath  string
		crc32h uint32
	)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fpath = path.Join(dir, file.Name())
		crc32h, err = crc32Hash(fpath)

		if err != nil {
			return nil, err
		}

		if _, ok := matches[crc32h]; ok {
			matches[crc32h] = append(matches[crc32h], fpath)
		} else {
			matches[crc32h] = []string{fpath}
		}
	}

	return matches, nil
}

func moveFiles(to string, oldpaths ...string) {
	var (
		newlocation string
		err         error
	)
	for _, op := range oldpaths {
		newlocation = path.Join(to, filepath.Base(op))
		err = os.Rename(op, newlocation)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var (
		dir, to string
		err     error
	)

	flag.StringVar(&dir, "dir", "", "Specifies the directory to search")
	flag.StringVar(&to, "to", "", "Specifies the directory to move duplicate files to")
	flag.Parse()

	fmt.Println(dir, to)

	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	duplicates, err := findDuplicates(dir)
	if err != nil {
		log.Fatal(err)
	}

	if to != "" {
		err = os.Mkdir(to, 0777)
		if err != nil {
			log.Fatal(err)
		}
		for _, dups := range duplicates {
			if len(dups) > 1 {
				moveFiles(to, dups...)
			}
		}
	}
}
