package main

import (
	"log"
	"sync"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type DownPath struct {
	remotePath string
	localPath  string
}

type Session struct {
	ssh    *ssh.Client
	sftp   *sftp.Client
	addr   string
	id     string
	passwd string
}

// count 만큼 sftp 처리용 go-roution을 만든다.
func CreateSessons(count int, addr string, id string, passwd string, ch <-chan *DownPath) {

	makedCount := 0
	var wait sync.WaitGroup

	for i := 0; i < count; i++ {

		log.Printf("create session %d", i)

		session := Session{
			addr:   addr,
			id:     id,
			passwd: passwd,
		}

		var err error
		session.sftp, session.ssh, err = SftpConnect(addr, id, passwd)
		defer func() {
			session.sftp.Close()
			session.ssh.Close()
		}()

		if err != nil {
			log.Println(err)
		} else {
			makedCount++
		}

		wait.Add(1)

		go func(session Session) {

			infinity := true

			for infinity {

				select {
				case file := <-ch:
					// 채널을 읽어서 다운받는다.
					err := SftpDown(session.sftp, file.remotePath, file.localPath)
					if err != nil {
						log.Printf("%v", err)
					}

				default:
					infinity = false
				}

			}

			wait.Done()

		}(session)
	}

	if makedCount == 0 {
		log.Println("no valid sftp connection.")
	}

	wait.Wait()
}
