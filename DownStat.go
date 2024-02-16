package main

import (
	"encoding/json"
	"fmt"
)

type DownFile struct {
	Path       string `json:"path"`
	RemoteSize int64  `json:"remoteSize"`
	LocalSize  int64  `json:"localSize"`
	Notify     string `json:"notify,omitempty"`
}

type DownStat struct {
	State     string     `json:"stat"`
	DownFiles []DownFile `json:"files"`
}

func (f *DownStat) ToJSON() ([]byte, error) {

	data, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("fail DownStat to json %v, %s", f, err)
	}

	return data, nil
}
