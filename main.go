package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var downInfo DownInfo

const (
	AppPort        = "8888"
	ConfigFileName = "default.yaml"
	ListFileName   = "files.lst"
)

func main() {

	downInfo = DownInfo{
		State: READY,
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

		log.Println("HandleFunc /")

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if downInfo.State == DOWNLOADING || downInfo.State == PREPAREDOWN {
			log.Println("alreay downloading. / ")
			w.Write([]byte(HtmlDownload()))
			return
		}

		// read yaml
		downInfo.LoadSftp(ConfigFileName)

		// read files list
		downInfo.LoadRemoteFiles(ListFileName)

		w.Write([]byte(HtmlRoot()))
	})

	// 다운로드 view
	http.HandleFunc("/download", func(w http.ResponseWriter, req *http.Request) {

		log.Println("HandleFunc /download")

		// 이미 다운 중이면, 다운 중인 화면을 보여준다.
		if downInfo.State == DOWNLOADING || downInfo.State == PREPAREDOWN {
			log.Println("alreay downloading. /download")
			w.Write([]byte(HtmlDownload()))
			return
		}

		err := req.ParseForm()
		if err != nil {
			log.Printf("req ParseForm %v\n", err)
		}

		// 설정 읽기
		downInfo = DownInfo{
			State:        READY,
			Ip:           req.PostFormValue("sftp-addr"),
			Port:         22,
			Id:           req.PostFormValue("sftp-id"),
			Password:     req.PostFormValue("sftp-pwd"),
			LocalDir:     req.PostFormValue("local-dir"),
			SessionCount: 10,
			OverWrite:    true,
			Files:        []File{},
		}

		downInfo.SetRemoteFiles(req.PostFormValue("file-list"))

		// 설정 저장
		downInfo.SaveSftp(ConfigFileName)

		downInfo.SaveRemoteFiles(ListFileName)

		w.Write([]byte(HtmlDownload()))

		go startDownload()
	})

	// 상태 전송을 위한
	http.HandleFunc("/downloading", func(w http.ResponseWriter, req *http.Request) {

		log.Println("HandleFunc /downloading")

		w.Header().Add("Content-Type", "application/json")

		buf := getState()

		log.Println(buf)

		w.Write([]byte(buf))
	})

	err := http.ListenAndServe(":"+AppPort, nil)
	if err != nil {
		log.Println(err)
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

// 원격 파일 상태 체크
func remoteFileCheck() {

	sft, ssh, err := SftpConnect(downInfo.Addr(), downInfo.Id, downInfo.Password)
	if err != nil {
		panic(err)
	}
	defer func() {
		sft.Close()
		ssh.Close()
	}()

	for i, file := range downInfo.Files {

		if file.Duplicate {
			continue
		}

		// ftp 서버의 파일
		date, size, err := SftpGetFileInfo(sft, file.Remote.Path)
		if err != nil {
			if strings.HasPrefix(err.Error(), "fail to connect") {
				panic(err)
			}
			log.Println(err)
		} else {

			downInfo.Files[i].Remote.Date = date
			downInfo.Files[i].Remote.Size = size
			downInfo.Files[i].Remote.IsExist = true
		}

		log.Println(downInfo.Files[i].Remote)
	}
}

func localFileCheck() {
	for i, file := range downInfo.Files {

		if !file.Remote.IsExist || file.Duplicate {
			continue
		}

		size, err := getFileSize(file.Local.Path)
		if err != nil {
			fmt.Printf("ERROR! %v", err)
		}

		// 다운받을 파일이 이미 있는 경우
		if size == file.Remote.Size && !downInfo.OverWrite {
			downInfo.Files[i].Skip = true
		}
	}
}

func startDownload() {

	log.Println("prepare downloand")
	downInfo.State = PREPAREDOWN

	// 중복 체크
	downInfo.DuplicateCheck()

	// 원격 파일 체크
	remoteFileCheck()

	// 다운받을 local dir 계산
	downInfo.SetLocalFiles()

	// 이미 다운 받은 파일이 있는지 체크
	localFileCheck()

	if downInfo.SessionCount == 0 || downInfo.SessionCount > downInfo.VaildFilesCount() {
		downInfo.SessionCount = downInfo.VaildFilesCount()
	}

	// 채널, 세션 풀 생성
	ch := make(chan *DownPath, 100)
	defer close(ch)

	// 대기열에 넣기 시작
	go func(ch chan<- *DownPath) {
		for _, file := range downInfo.Files {

			if !file.Remote.IsExist || file.Duplicate {
				continue
			}

			downInfo := DownPath{
				remotePath: file.Remote.Path,
				localPath:  file.Local.Path,
			}

			ch <- &downInfo
		}
	}(ch)

	log.Println("start downloand")
	downInfo.State = DOWNLOADING

	// 대기열 전송 시작
	CreateSessons(downInfo.SessionCount, downInfo.Addr(), downInfo.Id, downInfo.Password, ch)

	downInfo.State = DONE

	log.Println("end downloand")
}

func getState() string {

	stat := DownStat{}

	if downInfo.State == DONE {
		stat.State = "DONE"
	} else if downInfo.State == PREPAREDOWN {
		stat.State = "PREPARE DOWNLOAD"
	} else if downInfo.State == DOWNLOADING {
		stat.State = "DOWNLOADING"
	} else {
		stat.State = "READY"
	}

	for _, file := range downInfo.Files {

		downFile := DownFile{
			Path:       file.Remote.Path,
			RemoteSize: file.Remote.Size,
		}

		if downInfo.State == READY || downInfo.State == PREPAREDOWN {

			// 다운로드 전

			downFile.LocalSize = 0
			downFile.Notify = ""
		} else {
			// 다운로드 중
			if file.Skip {
				if !file.Remote.IsExist {
					downFile.Notify = "The system cannot find the file specified. remote"
				} else if file.Duplicate {
					downFile.Notify = "Duplicate file."
				}
			} else {
				size, err := getFileSize(file.Local.Path)
				if err != nil {
					downFile.LocalSize = 0
					downFile.Notify = err.Error()
				} else {
					downFile.LocalSize = size
					downFile.Notify = ""
				}
			}
		}

		stat.DownFiles = append(stat.DownFiles, downFile)
	}

	//	log.Printf("stat: %v\n", stat)

	buf, err := stat.ToJSON()
	if err != nil {
		log.Printf("fail stat to json %v, %s", stat, err)
	}

	return string(buf)
}

func getFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Printf("fail stat fileinfo %s, %s", path, err)
		return 0, err
	}

	return fi.Size(), nil
}
