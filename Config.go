package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Status int

const (
	READY Status = iota
	DOWNLOADING
	DONE
)

type Config struct {
	Ip          string     `yaml:"sftp.ip,omitempty"`
	Port        int        `yaml:"sftp.port,omitempty"`
	Id          string     `yaml:"sftp.id,omitempty"`
	Password    string     `yaml:"sftp.password,omitempty"`
	LocalDir    string     `yaml:"local.directory,omitempty"`
	RemoteFiles []FileInfo `yaml:"-"`
	LocalFiles  []FileInfo `yaml:"-"`
	State       Status     `yaml:"-"`
}

// Sftp 설정 파일을 읽는다.
func (conf *Config) LoadSftp(fileName string) error {

	buf, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read config file %s, ReadFile: %v", fileName, err)
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return fmt.Errorf("invaild config file %s, Unmarshal: %v", fileName, err)
	}

	if conf.Port == 0 {
		conf.Port = 22
	}

	return nil
}

// Sftp 설정 파일을 쓴다.
func (conf *Config) SaveSftp(fileName string) error {

	buf, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("fail marshal config %v, Marshal: %v", conf, err)
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("cannot open config file %s, WriteFile: %v", fileName, err)
	}
	defer file.Close()

	_, err = file.Write(buf)
	if err != nil {
		return fmt.Errorf("cannot write config file %s, WriteFile: %v", fileName, err)
	}

	return nil
}

// 목록 파일을 읽는다.
func (conf *Config) LoadRemoteFiles(fileName string) error {

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("fail file open %s, err: %v", fileName, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := strings.Trim(scanner.Text(), " \t\r")

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		conf.RemoteFiles = append(conf.RemoteFiles, FileInfo{path: line})
	}

	return nil
}

// 목록 파일을 쓴다.
func (conf *Config) SaveRemoteFiles(fileName string) error {

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("cannot open list file %s, WriteFile: %v", fileName, err)
	}
	defer file.Close()

	buf := ""
	for _, file := range conf.RemoteFiles {
		buf += fmt.Sprintln(file.path)
	}

	_, err = file.Write([]byte(buf))
	if err != nil {
		return fmt.Errorf("cannot write list file %s, WriteFile: %v", fileName, err)
	}

	return nil
}

// 개행 문자로 파일 목록을 파싱한다.
func (conf *Config) SetRemoteFiles(text string) {
	slice := strings.Split(text, "\n")
	for _, line := range slice {
		line = strings.Trim(line, " \t\r")

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		conf.RemoteFiles = append(conf.RemoteFiles, FileInfo{path: line})
	}
}

// LocalHome = LocalDir + RemoteFiles 공통 경로의 한단계 상위 경로,
// LocalFiles.path = LocalHome + 공통 경로를 제외한 RemoteFiles 경로
func (conf *Config) SetLocalFiles() {

	// 공통 경로
	common := GetParentDir(conf.RemoteFiles[0].path, "/")

	for _, file := range conf.RemoteFiles {
		common = GetSameDir(common, file.path)
	}

	for _, file := range conf.RemoteFiles {
		path := conf.LocalDir + "/" + strings.TrimPrefix(file.path, common)

		// local 디렉토리 구분자로 변경
		path = strings.ReplaceAll(path, "/", string(os.PathSeparator))

		conf.LocalFiles = append(conf.LocalFiles, FileInfo{path: path})
	}
}
