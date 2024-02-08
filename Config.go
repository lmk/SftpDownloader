package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sftp  Sftp  `yaml:"sftp.ip"`
	Local Local `yaml:"local"`
}

type Sftp struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Id       string `yaml:"id"`
	Password string `yaml:"password"`
}

type Local struct {
	Directory string `yaml:"directory"`
}

func (conf *Config) Load(fileName string) error {

	conf = &Config{Sftp{Ip: "", Port: 22, Id: "", Password: ""}, Local{Directory: ""}}

	buf, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read config file %s, ReadFile: %v", fileName, err)
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return fmt.Errorf("invaild config file %s, Unmarshal: %v", fileName, err)
	}

	return nil
}

func (conf *Config) Save(fileName string) error {

	buf, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("fail marshal config %v, Marshal: %v", conf, err)
	}

	err = os.WriteFile(fileName, buf, 0660)
	if err != nil {
		return fmt.Errorf("cannot write config file %s, WriteFile: %v", fileName, err)
	}

	return nil
}
