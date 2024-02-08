package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func SftpConn(config Config) (*sftp.Client, error) {

	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(config.Sftp.Password))

	sshConfig := ssh.ClientConfig{
		User: config.Sftp.Id,
		Auth: auths,
	}

	addr := fmt.Sprintf("%s:%d", config.Sftp.Ip, config.Sftp.Port)

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

func SftpDown(config Config, file *FileInfo) error {

	time.Sleep(10 * time.Second)

	return nil
}

func CheckFile(config Config, file *FileInfo) error {

	sc, err := SftpConn(config)
	if err != nil {
		return err
	}

	sc.ReadDir()

	return nil
}
