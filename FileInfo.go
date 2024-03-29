package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// FileInfo file 정보
type FileInfo struct {
	IsExist bool   `json:"exist,omitempty"`
	Size    int64  `json:"size,omitempty"`
	Date    string `json:"date,omitempty"`
	Path    string `json:"path,omitempty"`
}

// TheFileInfo 파일 정보를 읽는다.
func TheFileInfo(f string) FileInfo {

	var fi FileInfo

	fi.Path = f

	fs, err := os.Stat(f)

	fi.IsExist = !errors.Is(err, os.ErrNotExist)

	if fi.IsExist {
		fi.Size = fs.Size()
		fi.Date = fs.ModTime().Format("2006-01-02 15:04:05")
	}

	return fi
}

func (f *FileInfo) ExistToString() string {

	OX := "X"
	if f.IsExist {
		OX = "O"
	}

	//return fmt.Sprintf("%s %19s %7s %s", OX, f.date, HumanSize(float64(f.size)), f.path)
	return OX
}

func (f *FileInfo) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("fail fileInfo to json %v, %s", f, err)
	}

	return data, nil
}
