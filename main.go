package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var cfg Config

const (
	AppPort        = "8888"
	ConfigFileName = "default.yaml"
	ListFileName   = "files.lst"
)

func main() {

	cfg = Config{
		State:       READY,
		RemoteFiles: []FileInfo{},
		LocalFiles:  []FileInfo{},
	}

	go runServer()

	waitListen(AppPort)

	openBrowser("http://localhost:" + AppPort + "/")

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

// 웹브라우져를 열어서 UI를 띄운다.
func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func runServer() {

	// 첫 페이지
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		fmt.Println("HandleFunc /")

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if cfg.State == DOWNLOADING {
			fmt.Println("alreay downloading. / ")
			w.Write([]byte(HtmlDownload()))
			return
		}

		// read yaml
		cfg.LoadSftp(ConfigFileName)

		// read files list
		cfg.LoadRemoteFiles(ListFileName)

		w.Write([]byte(HtmlRoot()))
	})

	// 다운로드 view
	http.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {

		fmt.Println("HandleFunc /download")

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if cfg.State == DOWNLOADING {
			fmt.Println("alreay downloading. /download")
			w.Write([]byte(HtmlDownload()))
			return
		}

		err := req.ParseForm()
		if err != nil {
			fmt.Printf("req ParseForm %v\n", err)
		}

		// 설정 읽기
		cfg = Config{
			Ip:          req.PostFormValue("sftp-addr"),
			Port:        22,
			Id:          req.PostFormValue("sftp-id"),
			Password:    req.PostFormValue("sftp-pwd"),
			LocalDir:    req.PostFormValue("local-dir"),
			RemoteFiles: []FileInfo{},
			State:       READY,
		}

		cfg.SetRemoteFiles(req.PostFormValue("file-list"))

		// 설정 저장
		cfg.SaveSftp(ConfigFileName)

		cfg.SaveRemoteFiles(ListFileName)

		// local dir 계산
		cfg.SetLocalFiles()

		w.Write([]byte(HtmlDownload()))

		go startDownload()
	})

	// 상태 전송을 위한
	http.HandleFunc("/downloading", func(w http.ResponseWriter, req *http.Request) {

		fmt.Println("HandleFunc /downloading")

		w.Header().Add("Content-Type", "application/json")

		w.Write([]byte(getState()))
	})

	err := http.ListenAndServe(":"+AppPort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// port 가 listen 상태가 될때까지 대기한다
func waitListen(port string) bool {
	conn, err := net.DialTimeout("tcp", "localhost:"+port, 60*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func startDownload() {

	cfg.State = DOWNLOADING

	var wait sync.WaitGroup

	for i, file := range cfg.RemoteFiles {

		// ftp 서버의 파일
		date, size, err := SftpGetFileInfo(cfg, file.path)
		if err != nil {
			if strings.HasPrefix(err.Error(), "fail to connect") {
				panic(err)
			}
			fmt.Println(err)
		} else {

			cfg.RemoteFiles[i].date = date
			cfg.RemoteFiles[i].size = size

			wait.Add(1)
			// download
			go func(i int) {
				SftpDown(cfg, cfg.RemoteFiles[i].path, cfg.LocalFiles[i].path)
				wait.Done()
			}(i)
		}
	}

	wait.Wait()

	cfg.State = DONE
}

func getState() string {

	stat := DownStat{}

	if cfg.State == DONE {
		stat.State = "DONE"
	} else if cfg.State == DOWNLOADING {
		stat.State = "DOWNLOADING"
	} else {
		stat.State = "READY"
	}

	for i := range cfg.RemoteFiles {
		downFile := DownFile{
			path:       cfg.RemoteFiles[i].path,
			remoteSize: cfg.RemoteFiles[i].size,
			localSize:  getFileSize(cfg.LocalFiles[i].path),
		}

		stat.DownFiles = append(stat.DownFiles, downFile)
	}

	buf, err := stat.MarshalJSON()
	if err != nil {
		fmt.Printf("fail stat to json %v, %s", stat, err)
	}

	return string(buf)
}

func getFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Printf("fail stat fileinfo %s, %s", path, err)
	}

	return fi.Size()
}
