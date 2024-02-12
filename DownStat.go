package main

import (
	"encoding/json"
	"fmt"
)

type DownStat struct {
	State     string     `json:"stat"`
	DownFiles []DownFile `json:"files"`
}

type DownFile struct {
	path       string `json:"path"`
	remoteSize int64  `json:"remoteSize"`
	localSize  int64  `json:"localSize"`
}

func (f *DownStat) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("fail DownStat to json %v, %s", f, err)
	}

	return data, nil
}
