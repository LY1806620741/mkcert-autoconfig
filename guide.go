package main

type Guide struct {
	prompt *prompt
}

func guideRun() {
	(&Guide{&prompt{}}).Run()
}

func (g *Guide) Run() {
	//检查是否已经拥有根证书
	g.prompt.GenRootCert()
}
