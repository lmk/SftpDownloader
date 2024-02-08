package main

import (
	"encoding/json"
	"fmt"
)

type Status int

const (
	ready Status = iota
	downloading
	done
)

type State struct {
	State Status     `json:"state,omitempty"`
	Files []FileInfo `json:"files,omitempty"`
}

func (f *State) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("fail state to json %v, %s", f, err)
	}

	return data, nil
}
