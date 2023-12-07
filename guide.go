package main
import(
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
	m:=mkcert{}
	m.CAROOT="./"
	newCA(m)
	if (false){
		if (g.prompt.GenRootCert()){
			
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

//初始化ca
func newCA(m mkcert){
}

// //
// func addFile(em){
	
// }
