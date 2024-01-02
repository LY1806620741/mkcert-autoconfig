// Copyright 2018 The mkcert Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command mkcert is a simple zero-config tool to make development certificates.
package main

import (
	"crypto"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"golang.org/x/net/idna"
)

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var Version string

// 独立开，不干扰原来mkcert的逻辑，明确这个是一个mkcert的向导工具
func main() {
	if len(os.Args) == 1 {
		fmt.Print(i18nText.help)
		return
	}
	log.SetFlags(0)
	flag.Parse()
	//检查是否在命令列表
	mainCommand := flag.Arg(0)

	switch mainCommand {
	case "auto":
		guideRun()
		break
	case "mkcert":
		//删除第二个参数
		os.Args = append(os.Args[:1], os.Args[2:]...)
		//运行原mkcert逻辑
		mkcertMain()
		break
	default:
		log.Fatalln(i18nText.errUnknownGroupCommand(mainCommand))
	}
}

func mkcertMain() {

	if len(os.Args) == 1 {
		fmt.Print(i18nMkcertText.shortUsage)
		return
	}
	log.SetFlags(0)
	var (
		installFlag	= flag.Bool("install", false, "")
		uninstallFlag	= flag.Bool("uninstall", false, "")
		pkcs12Flag	= flag.Bool("pkcs12", false, "")
		ecdsaFlag	= flag.Bool("ecdsa", false, "")
		clientFlag	= flag.Bool("client", false, "")
		helpFlag	= flag.Bool("help", false, "")
		carootFlag	= flag.Bool("CAROOT", false, "")
		csrFlag		= flag.String("csr", "", "")
		certFileFlag	= flag.String("cert-file", "", "")
		keyFileFlag	= flag.String("key-file", "", "")
		p12FileFlag	= flag.String("p12-file", "", "")
		versionFlag	= flag.Bool("version", false, "")
	)
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), i18nMkcertText.shortUsage)
		fmt.Fprintln(flag.CommandLine.Output(), i18nMkcertText.moreOptions)
	}
	flag.Parse()
	if *helpFlag {
		fmt.Print(i18nMkcertText.shortUsage)
		fmt.Print(i18nMkcertText.advancedUsage)
		return
	}
	if *versionFlag {
		if Version != "" {
			fmt.Println(Version)
			return
		}
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}
	if *carootFlag {
		if *installFlag || *uninstallFlag {
			log.Fatalln(i18nText.scan50,
			)
		}
		fmt.Println(getCAROOT())
		return
	}
	if *installFlag && *uninstallFlag {
		log.Fatalln(i18nText.scan51,
		)
	}
	if *csrFlag != "" && (*pkcs12Flag || *ecdsaFlag || *clientFlag) {
		log.Fatalln(i18nText.scan52,
		)
	}
	if *csrFlag != "" && flag.NArg() != 0 {
		log.Fatalln(i18nText.scan53,
		)
	}
	(&mkcert{
		installMode:	*installFlag, uninstallMode: *uninstallFlag, csrPath: *csrFlag,
		pkcs12:	*pkcs12Flag, ecdsa: *ecdsaFlag, client: *clientFlag,
		certFile:	*certFileFlag, keyFile: *keyFileFlag, p12File: *p12FileFlag,
	}).Run(flag.Args())
}

const rootName = "rootCA.pem"
const rootKeyName = "rootCA-key.pem"

type mkcert struct {
	installMode, uninstallMode	bool
	pkcs12, ecdsa, client		bool
	keyFile, certFile, p12File	string
	csrPath				string

	CAROOT	string
	caCert	*x509.Certificate
	caKey	crypto.PrivateKey

	// The system cert pool is only loaded once. After installing the root, checks
	// will keep failing until the next execution. TODO: maybe execve?
	// https://github.com/golang/go/issues/24540 (thanks, myself)
	ignoreCheckFailure	bool
}

func (m *mkcert) Run(args []string) {
	m.CAROOT = getCAROOT()
	if m.CAROOT == "" {
		log.Fatalln(i18nText.scan54,
		)
	}
	fatalIfErr(os.MkdirAll(m.CAROOT, 0755), i18nText.scan55,
	)
	m.loadCA()

	if m.installMode {
		m.install()
		if len(args) == 0 {
			return
		}
	} else if m.uninstallMode {
		m.uninstall()
		return
	} else {
		var warning bool
		if storeEnabled("system") && !m.checkPlatform() {
			warning = true
			log.Println(i18nText.scan56,
			)
		}
		if storeEnabled("nss") && hasNSS && CertutilInstallHelp != "" && !m.checkNSS() {
			warning = true
			log.Printf(i18nText.scan57,

				NSSBrowsers)
		}
		if storeEnabled("java") && hasJava && !m.checkJava() {
			warning = true
			log.Println(i18nText.scan58,
			)
		}
		if warning {
			log.Println(i18nText.scan59,
			)
		}
	}

	if m.csrPath != "" {
		m.makeCertFromCSR()
		return
	}

	if len(args) == 0 {
		flag.Usage()
		return
	}

	hostnameRegexp := regexp.MustCompile(`(?i)^(\*\.)?[0-9a-z_-]([0-9a-z._-]*[0-9a-z_-])?$`)
	for i, name := range args {
		if ip := net.ParseIP(name); ip != nil {
			continue
		}
		if email, err := mail.ParseAddress(name); err == nil && email.Address == name {
			continue
		}
		if uriName, err := url.Parse(name); err == nil && uriName.Scheme != "" && uriName.Host != "" {
			continue
		}
		punycode, err := idna.ToASCII(name)
		if err != nil {
			log.Fatalf(i18nText.scan60,

				name, err)
		}
		args[i] = punycode
		if !hostnameRegexp.MatchString(punycode) {
			log.Fatalf(i18nText.scan61,

				name)
		}
	}

	m.makeCert(args)
}

func getCAROOT() string {
	if env := os.Getenv("CAROOT"); env != "" {
		return env
	}

	var dir string
	switch {
	case runtime.GOOS == "windows":
		dir = os.Getenv("LocalAppData")
	case os.Getenv("XDG_DATA_HOME") != "":
		dir = os.Getenv("XDG_DATA_HOME")
	case runtime.GOOS == "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			return ""
		}
		dir = filepath.Join(dir, "Library", "Application Support")
	default:	// Unix
		dir = os.Getenv("HOME")
		if dir == "" {
			return ""
		}
		dir = filepath.Join(dir, ".local", "share")
	}
	return filepath.Join(dir, "mkcert")
}

func (m *mkcert) install() {
	if storeEnabled("system") {
		if m.checkPlatform() {
			log.Print(i18nText.scan48,
			)
		} else {
			if m.installPlatform() {
				log.Print(i18nText.scan62,
				)
			}
			m.ignoreCheckFailure = true	// TODO: replace with a check for a successful install
		}
	}
	if storeEnabled("nss") && hasNSS {
		if m.checkNSS() {
			log.Printf(i18nText.scan63,

				NSSBrowsers)
		} else {
			if hasCertutil && m.installNSS() {
				log.Printf(i18nText.scan64,

					NSSBrowsers)
			} else if CertutilInstallHelp == "" {
				log.Printf(i18nText.scan65,

					NSSBrowsers)
			} else if !hasCertutil {
				log.Printf(i18nText.scan66,

					NSSBrowsers)
				log.Printf(i18nText.scan67,

					CertutilInstallHelp)
			}
		}
	}
	if storeEnabled("java") && hasJava {
		if m.checkJava() {
			log.Println(i18nText.scan68,
			)
		} else {
			if hasKeytool {
				m.installJava()
				log.Println(i18nText.scan69,
				)
			} else {
				log.Println(i18nText.scan70,
				)
			}
		}
	}
	log.Print(i18nText.scan71,
	)
}

func (m *mkcert) uninstall() {
	if storeEnabled("nss") && hasNSS {
		if hasCertutil {
			m.uninstallNSS()
		} else if CertutilInstallHelp != "" {
			log.Print(i18nText.scan71,
			)
			log.Printf(i18nText.scan72,

				NSSBrowsers)
			log.Printf(i18nText.scan73,

				CertutilInstallHelp)
			log.Print(i18nText.scan71,
			)
		}
	}
	if storeEnabled("java") && hasJava {
		if hasKeytool {
			m.uninstallJava()
		} else {
			log.Print(i18nText.scan71,
			)
			log.Println(i18nText.scan74,
			)
			log.Print(i18nText.scan71,
			)
		}
	}
	if storeEnabled("system") && m.uninstallPlatform() {
		log.Print(i18nText.scan75,
		)
		log.Print(i18nText.scan71,
		)
	} else if storeEnabled("nss") && hasCertutil {
		log.Printf(i18nText.scan76,

			NSSBrowsers)
		log.Print(i18nText.scan71,
		)
	}
}

func (m *mkcert) checkPlatform() bool {
	if m.ignoreCheckFailure {
		return true
	}

	_, err := m.caCert.Verify(x509.VerifyOptions{})
	return err == nil
}

func storeEnabled(name string) bool {
	stores := os.Getenv("TRUST_STORES")
	if stores == "" {
		return true
	}
	for _, store := range strings.Split(stores, ",") {
		if store == name {
			return true
		}
	}
	return false
}

func fatalIfErr(err error, msg string) {
	if err != nil {
		log.Fatalf(i18nText.scan77,

			msg, err)
	}
}

func fatalIfCmdErr(err error, cmd string, out []byte) {
	if err != nil {
		log.Fatalf(i18nText.scan78,

			cmd, err, out)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func binaryExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

var sudoWarningOnce sync.Once

func commandWithSudo(cmd ...string) *exec.Cmd {
	if u, err := user.Current(); err == nil && u.Uid == "0" {
		return exec.Command(cmd[0], cmd[1:]...)
	}
	if !binaryExists("sudo") {
		sudoWarningOnce.Do(func() {
			log.Println(i18nText.scan79,
			)
		})
		return exec.Command(cmd[0], cmd[1:]...)
	}
	return exec.Command("sudo", append([]string{"--prompt=Sudo password:", "--"}, cmd...)...)
}
