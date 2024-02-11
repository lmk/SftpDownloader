package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
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

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if cfg.State == DOWNLOADING {
			fmt.Println("alreay downloading.")
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

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if cfg.State == DOWNLOADING {
			fmt.Println("alreay downloading.")
			w.Write([]byte(HtmlDownload()))
			return
		}

		err := req.ParseForm()
		if err != nil {
			fmt.Printf("req ParseForm %v\n", err)
		}

		fmt.Println(req.Form)

		// 설정 읽기
		cfg = Config{
			Ip:          req.PostFormValue("sftp-addr"),
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

	for i := range cfg.RemoteFiles {
		wait.Add(1)

		// ftp 서버의 파일
		err := SftpGetFileInfo(cfg, &cfg.RemoteFiles[i])
		if err != nil {
			fmt.Println(err)
		}

		go SftpDown(cfg, cfg.RemoteFiles[i].path, cfg.LocalFiles[i].path)
	}

	wait.Wait()

	cfg.State = DONE
}

func getState() string {

	for _, file := range fileList.Files {

	}

	return ""
}
