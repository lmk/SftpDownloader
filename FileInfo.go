package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// FileInfo file 정보
type FileInfo struct {
	isExists bool   `json:"exist,omitempty"`
	size     int64  `json:"size,omitempty"`
	date     string `json:"date,omitempty"`
	path     string `json:"path,omitempty"`
	doing    int    `json:"doing,omitempty"` // 다운로드 진행률
}

// TheFileInfo 파일 정보를 읽는다.
func TheFileInfo(f string) FileInfo {

	var fi FileInfo

	fi.path = f

	fs, err := os.Stat(f)

	fi.isExists = !errors.Is(err, os.ErrNotExist)

	if fi.isExists {
		fi.size = fs.Size()
		fi.date = fs.ModTime().Format("2006-01-02 15:04:05")
	}

	return fi
}

func (f *FileInfo) ToString() string {

	OX := "X"
	if f.isExists {
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
