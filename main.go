package main

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var config Config
var fileList FileList

var state State

const (
	AppPort        = "8888"
	ConfigFileName = "default.yaml"
	ListFileName   = "files.lst"
)

func main() {

	state = State{State: ready}

	go runServer()

	waitListen(AppPort)

	openBrowser("http://localhost:" + AppPort + "/")

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

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

		// read yaml
		config.Load(ConfigFileName)

		// read pkg list
		fileList.Load(ListFileName)

		w.Write([]byte(HtmlRoot(config, fileList)))
	})

	// 다운로드 view
	http.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {

		if state.State == downloading {
			fmt.Println("alreay downloading.")
			return
		}

		config = Config{
			Sftp: Sftp{
				Ip:       req.PostFormValue("sftp-addr"),
				Id:       req.PostFormValue("sftp-id"),
				Password: req.PostFormValue("sftp-pwd"),
			},
			Local: Local{
				Directory: req.PostFormValue("local-dir"),
			},
		}

		fileList = FileList{
			Files: []FileInfo{},
		}

		fileList.FromString(req.PostFormValue("file-list"))

		config.Save(ConfigFileName)

		fileList.Save(ListFileName)

		w.Write([]byte(HtmlDownload(fileList)))
	})

	// 상태 전송을 위한
	http.HandleFunc("/downloading", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		w.Write([]byte(GetState()))
	})

	err := http.ListenAndServe(":"+AppPort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func waitListen(port string) bool {
	conn, err := net.DialTimeout("tcp", "localhost:"+port, 60*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func GetStat() {

}
