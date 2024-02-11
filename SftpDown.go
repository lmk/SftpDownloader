package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func SftpConn(config Config) (*sftp.Client, error) {

	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(config.Password))

	sshConfig := ssh.ClientConfig{
		User: config.Id,
		Auth: auths,
	}

	addr := fmt.Sprintf("%s:%d", config.Ip, config.Port)

	conn, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connecto to [%s]: %v\n", addr, err)
	}
	defer conn.Close()

	sc, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("Unable to start SFTP subsystem: %v\n", err)
	}
	defer sc.Close()

	return sc, nil
}

func SftpDown(config Config, remote string, local string) error {

	fmt.Printf("download %s to %s\n", remote, local)
	time.Sleep(10 * time.Second)

	return nil
}

func SftpGetFileInfo(config Config, file *FileInfo) error {

	sc, err := SftpConn(config)
	if err != nil {
		return fmt.Errorf("fail CheckFile %v, %v", config, err)
	}

	fi, err := sc.Stat(file.path)
	if err != nil {
		return fmt.Errorf("fail CheckFile %v, %v", file.path, err)
	}

	file.date = fi.ModTime().Format("2006-01-02 15:04:05")
	file.size = fi.Size()
	file.isExists = true

	return nil
}
