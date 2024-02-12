package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func SftpDown(config Config, remote string, local string) error {

	fmt.Printf("download %s to %s\n", remote, local)
	time.Sleep(10 * time.Second)

	return nil
}

// get sftp remote file info
func SftpGetFileInfo(config Config, path string) (string, int64, error) {

	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(config.Password))

	addr := fmt.Sprintf("%s:%d", config.Ip, config.Port)

	sshConfig := ssh.ClientConfig{
		User:            config.Id,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		return "", 0, fmt.Errorf("fail to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		return "", 0, fmt.Errorf("unable to start SFTP subsystem: %v", err)
	}
	defer sc.Close()

	fi, err := sc.Stat(path)
	if err != nil {
		return "", 0, fmt.Errorf("fail SftpGetFileInfo %v, %v", path, err)
	}

	return fi.ModTime().Format("2006-01-02 15:04:05"),
		fi.Size(),
		nil
}
