package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

func SftpDown(sc *sftp.Client, remote string, local string) error {

	log.Printf("download %s to %s\n", remote, local)

	// local dir 생성
	dir := GetParentDir(local, string(os.PathSeparator))
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("fail local mkdir: %v", err)
	}

	// download file
	srcFile, err := sc.OpenFile(remote, (os.O_RDONLY))
	if err != nil {
		return fmt.Errorf("fail open remote: %v", err)
	}

	dstFile, err := os.Create(local)
	if err != nil {
		return fmt.Errorf("fail create local: %v", err)
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("fail copy : %v", err)
	}

	return nil
}

// get sftp remote file info
func SftpGetFileInfo(sc *sftp.Client, path string) (string, int64, error) {

	// get info
	fi, err := sc.Stat(path)
	if err != nil {
		return "", 0, fmt.Errorf("fail SftpGetFileInfo %v, %v", path, err)
	}

	return fi.ModTime().Format("2006-01-02 15:04:05"),
		fi.Size(),
		nil
}

func SftpConnect(addr, id, passwd string) (*sftp.Client, *ssh.Client, error) {

	// connect ssh
	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(passwd))

	sshConfig := ssh.ClientConfig{
		User:            id,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to connect to [%s]: %v", addr, err)
	}

	sc, err := sftp.NewClient(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to start SFTP subsystem: %v", err)
	}

	return sc, conn, nil
}
