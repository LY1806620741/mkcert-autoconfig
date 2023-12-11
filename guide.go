package main

import (
	"fmt"
	"log"
)

type Guide struct {
	prompt *prompt
}

func guideRun() {
	(&Guide{&prompt{}}).Run()
}

func (g *Guide) Run() {
	//检查是否已经拥有根证书
	m := mkcert{}
	m.CAROOT = "./"
	if caInit {
		text, _ := selffs.ReadFileString("key.pem")
		fmt.Println(text)
	} else {
		newCA(m)
		division()
	}
	if false {
		if g.prompt.GenRootCert() {

			if m.checkPlatform() {
				log.Print("The local CA is already installed in the system trust store! 👍")
			} else {
				// if m.installPlatform() {
				// 	log.Print("The local CA is now installed in the system trust store! ⚡️")
				// }
				m.ignoreCheckFailure = true // TODO: replace with a check for a successful install
			}
		}
	}
}

// 初始化ca
func newCA(m mkcert) {
	selffs.WriteFile("key.pem", "thisji")
	text, _ := selffs.ReadFileString("key.pem")
	fmt.Println(text)
}
