package main

import (
	"io/ioutil"
	"log"
	"os/user"
	"path"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestSSHClient(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		panic(u)
	}
	key, err := ioutil.ReadFile(path.Join(u.HomeDir, ".ssh", "id_rsa"))
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		//Auth: []ssh.AuthMethod{
		//	ssh.Password("dev2020!"),
		//},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "192.168.1.41:22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	buf, err := session.CombinedOutput("df -h")
	t.Logf("Command result: %s", string(buf))
	// sftpClient, err := sftp.NewClient(client)
	// if err != nil { //创建客户端
	// 	fmt.Println("创建客户端失败", err)
	// 	return
	// }
	// dstFile, err := sftpClient.Create("/root/tabby-1.0.168-setup.exe")
	// writer := bufio.NewWriter(dstFile)
	// defer func() {
	// 	writer.Flush()
	// 	dstFile.Close()
	// }()
	// localPath := "C:\\Users\\eices\\Downloads\\tabby-1.0.168-setup.exe"
	// srcFile, err := os.Open(localPath)
	// if err != nil {
	// 	fmt.Println("打开文件失败", err)
	// 	panic(err)
	// }
	// defer srcFile.Close()
	// reader := bufio.NewReader(srcFile)

	// written, err := io.Copy(writer, reader)
	// if err != nil {
	// 	t.Errorf("%s", err.Error())
	// }

	// t.Logf("Upload file is size: %d", written)
}
