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
	PREPAREDOWN
	DOWNLOADING
	DONE
)

type File struct {
	Remote    FileInfo
	Local     FileInfo
	Duplicate bool
	Skip      bool
}

type DownInfo struct {
	Ip             string `yaml:"sftp.ip,omitempty"`
	Port           int    `yaml:"sftp.port,omitempty"`
	Id             string `yaml:"sftp.id,omitempty"`
	Password       string `yaml:"sftp.password,omitempty"`
	LocalDir       string `yaml:"local.directory,omitempty"`
	LocalDirOption string `yaml:"local.dir-option,omitempty"`
	SessionCount   int    `yaml:"sftp.session-count,omitempty"`
	OverWrite      *bool  `yaml:"sftp.over-write,omitempty"`
	State          Status `yaml:"-"`
	Files          []File `yaml:"-"`
}

// Sftp 설정 파일을 읽는다.
func (info *DownInfo) LoadSftp(fileName string) error {

	buf, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read config file %s, ReadFile: %v", fileName, err)
	}

	err = yaml.Unmarshal(buf, info)
	if err != nil {
		return fmt.Errorf("invaild config file %s, Unmarshal: %v", fileName, err)
	}

	if info.Port == 0 {
		info.Port = 22
	}

	return nil
}

// Sftp 설정 파일을 쓴다.
func (info *DownInfo) SaveSftp(fileName string) error {

	buf, err := yaml.Marshal(info)
	if err != nil {
		return fmt.Errorf("fail marshal config %v, Marshal: %v", info, err)
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
func (info *DownInfo) LoadRemoteFiles(fileName string) error {

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

		info.Files = append(info.Files, File{Remote: FileInfo{Path: line}})
	}

	return nil
}

// 목록 파일을 쓴다.
func (info *DownInfo) SaveRemoteFiles(fileName string) error {

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("cannot open list file %s, WriteFile: %v", fileName, err)
	}
	defer file.Close()

	buf := ""
	for _, file := range info.Files {
		buf += fmt.Sprintln(file.Remote.Path)
	}

	_, err = file.Write([]byte(buf))
	if err != nil {
		return fmt.Errorf("cannot write list file %s, WriteFile: %v", fileName, err)
	}

	return nil
}

// 개행 문자로 파일 목록을 파싱한다.
func (info *DownInfo) SetRemoteFiles(text string) {
	slice := strings.Split(text, "\n")
	for _, line := range slice {
		line = strings.Trim(line, " \t\r")

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// 중복 체크
		for i, file := range info.Files {
			if file.Remote.Path == line {
				info.Files[i].Duplicate = true
				break
			}
		}

		info.Files = append(info.Files, File{Remote: FileInfo{Path: line}})
	}
}

// LocalHome = LocalDir + RemoteFiles 공통 경로의 한단계 상위 경로,
// LocalFiles.path = LocalHome + 공통 경로를 제외한 RemoteFiles 경로
func (info *DownInfo) SetLocalFiles() {

	// 공통 경로
	common := GetParentDir(info.Files[0].Remote.Path, "/")

	for _, file := range info.Files {
		common = GetSameDir(common, file.Remote.Path)
	}

	for i, file := range info.Files {

		path := ""
		// switch info.LocalDirOption {
		// case "same-local": // retmoe 경로에서 LocalDir 과 동일한 경로를 찾아서 이하 dir을 생성하여 각 dir에 다운로드
		// 	path = info.LocalDir + "/" + strings.TrimPrefix(file.Remote.Path, common)
		// case "smart": // remote에서 공통 경로를 찾아서 localDir에 공통 경로 이하 dir을 생성하여 각 dir에 다운로드
		// default:
		path = info.LocalDir + "/" + strings.TrimPrefix(file.Remote.Path, common)
		// }

		// local 디렉토리 구분자로 변경
		info.Files[i].Local.Path = strings.ReplaceAll(path, "/", string(os.PathSeparator))
	}
}

// 다운받을 유효한 파일 개수
func (info *DownInfo) VaildFilesCount() int {

	count := 0

	for _, file := range info.Files {
		if file.Remote.IsExist && !file.Duplicate && !file.Skip {
			count++
		}
	}

	return count
}

func (info *DownInfo) Addr() string {
	return fmt.Sprintf("%s:%d", info.Ip, info.Port)
}

// remote 파일 경로를 multi-line string으로 반환한다.
func (info *DownInfo) RemoteFilesString() string {
	buf := ""
	for _, file := range info.Files {
		buf += fmt.Sprintln(file.Remote.Path)
	}

	return buf
}

// 원격 파일 중복 체크
func (info *DownInfo) DuplicateCheck() {
	for i := range info.Files {
		for j := 0; j < i; j++ {
			if info.Files[i] == info.Files[j] {
				info.Files[i].Duplicate = true
			}
		}
	}
}
