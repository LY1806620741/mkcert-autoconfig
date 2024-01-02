package main

import (
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

const rootName = "rootCA.pem"

type mkcert struct {
	CAROOT	string
	caCert	*x509.Certificate

	// The system cert pool is only loaded once. After installing the root, checks
	// will keep failing until the next execution. TODO: maybe execve?
	// https://github.com/golang/go/issues/24540 (thanks, myself)
	ignoreCheckFailure	bool
}

func (m *mkcert) checkPlatform() bool {
	if m.ignoreCheckFailure {
		return true
	}

	_, err := m.caCert.Verify(x509.VerifyOptions{})
	return err == nil
}

func (m *mkcert) install() {
	if storeEnabled("system") {
		if m.checkPlatform() {
			log.Print(i18nText.
				scan95,
			)
		} else {
			if m.installPlatform() {
				log.Print(i18nText.
					scan95,
				)
			}
			m.ignoreCheckFailure = true	// TODO: replace with a check for a successful install
		}
	}
	if storeEnabled("nss") && hasNSS {
		if m.checkNSS() {
			log.Printf(i18nText.
				scan95,

				NSSBrowsers)
		} else {
			if hasCertutil && m.installNSS() {
				log.Printf(i18nText.
					scan95,

					NSSBrowsers)
			} else if CertutilInstallHelp == "" {
				log.Printf(i18nText.
					scan65,

					NSSBrowsers)
			} else if !hasCertutil {
				log.Printf(i18nText.
					scan66,

					NSSBrowsers)
				log.Printf(i18nText.
					scan67,

					CertutilInstallHelp)
			}
		}
	}
	if storeEnabled("java") && hasJava {
		if m.checkJava() {
			log.Println(i18nText.
				scan95,
			)
		} else {
			if hasKeytool {
				m.installJava()
				log.Println(i18nText.
					scan95,
				)
			} else {
				log.Println(i18nText.
					scan70,
				)
			}
		}
	}
	log.Print(i18nText.
		scan95,
	)
}

func (m *mkcert) uninstall() {
	if storeEnabled("nss") && hasNSS {
		if hasCertutil {
			m.uninstallNSS()
		} else if CertutilInstallHelp != "" {
			log.Print(i18nText.
				scan95,
			)
			log.Printf(i18nText.
				scan72,

				NSSBrowsers)
			log.Printf(i18nText.
				scan73,

				CertutilInstallHelp)
			log.Print(i18nText.
				scan95,
			)
		}
	}
	if storeEnabled("java") && hasJava {
		if hasKeytool {
			m.uninstallJava()
		} else {
			log.Print(i18nText.
				scan95,
			)
			log.Println(i18nText.
				scan74,
			)
			log.Print(i18nText.
				scan95,
			)
		}
	}
	if storeEnabled("system") && m.uninstallPlatform() {
		log.Print(i18nText.
			scan95,
		)
		log.Print(i18nText.
			scan95,
		)
	} else if storeEnabled("nss") && hasCertutil {
		log.Printf(i18nText.
			scan95,

			NSSBrowsers)
		log.Print(i18nText.
			scan95,
		)
	}
}

type item struct {
	Name		string
	Description	string
}

func (m *mkcert) caUniqueName() string {
	return "mkcert development CA " + m.caCert.SerialNumber.String()
}

func main() {

	if !caInit {
		panic(errors.New("异常，没有根证书"))
	}

	os.WriteFile(rootName, cert, os.ModePerm)

	items := []item{
		{Name: "安装根证书", Description: "当前具备根证书公钥，可以选择信任该证书签发的网站"},
		{Name: "卸载根证书", Description: "不再信任该证书签发的网站"},
		{Name: "退出", Description: "什么都不做"},
	}

	templates := &promptui.SelectTemplates{
		Label:		"{{ . }}?",
		Active:		"\U0001F336 {{ .Name | cyan }}",
		Inactive:	"  {{ .Name | cyan }}",
		Selected:	"\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- 详情 ----------
{{ "名字:" | faint }}	{{ .Name }}
{{ "解释:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := items[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:		"当前是授信客户端,你要做什么",
		Items:		items,
		Templates:	templates,
		Size:		4,
		Searcher:	searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		log.Fatalln(i18nText.
			scan95,
		)
	} else {
		m := &mkcert{}
		m.CAROOT = "./"
		certDERBlock, _ := pem.Decode(cert)
		if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
			log.Fatalln(i18nText.
				scan95,
			)
		}
		m.caCert, err = x509.ParseCertificate(certDERBlock.Bytes)
		fatalIfErr(err, i18nText.
			scan95,
		)
		if i == 0 {
			m.install()
		} else if i == 1 {
			m.uninstall()
		}
	}
}
