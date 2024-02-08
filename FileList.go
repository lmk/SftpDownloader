package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FileList struct {
	Files []FileInfo
}

func (list *FileList) Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("fail file open %s, err: %v", fileName, err)
	}
	defer f.Close()

	buf := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buf += scanner.Text()
	}

	list.FromString(buf)

	return nil
}

func (list *FileList) ToString() string {
	buf := ""

	for _, file := range list.Files {
		buf += fmt.Sprintln(file.path)
	}

	return buf
}

// list의 FileInfo.path 값을 채운다.
func (list *FileList) FromString(text string) error {

	list.Files = []FileInfo{}

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.Trim(line, " \t")

		// 주석 & 빈줄 처리
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		list.Files = append(list.Files, FileInfo{path: line})
	}

	return nil
}

func (list *FileList) Save(fileName string) error {

	buf := ""
	for _, file := range list.Files {
		buf += fmt.Sprintln("%s", file.path)
	}

	err := os.WriteFile(fileName, []byte(buf), 0660)
	if err != nil {
		return fmt.Errorf("cannot write list file %s, WriteFile: %v", fileName, err)
	}

	return nil
}
