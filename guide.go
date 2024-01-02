package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"log"
	"os"
	"strings"
	"time"
)

type Guide struct {
	prompt Prompt
}

func guideRun() {
	(&Guide{&prompt{}}).Run()
}

func (g *Guide) Run() {
	//检查是否已经拥有根证书
	m := mkcert{}
	m.CAROOT = "./"

	if caInit {
		loadCA(&m)
		// 母客户端菜单
		rootIndex := g.prompt.RootMenu()

		//生成子证书
		if rootIndex == 0 {
			m.makeCert(g.prompt.InputHost())
		} else if rootIndex == 1 {	//导出根证书
			keyContent, _ := selffs.ReadFile(rootKeyName)
			os.WriteFile(rootKeyName, keyContent, 0666)
			pemContent, _ := selffs.ReadFile(rootName)
			os.WriteFile(rootName, pemContent, 0666)
			log.Println(i18nText.
				scan46,
			)
		} else if rootIndex == 2 {	//导出客户端
			//生成html或客户端
			//解压客户端

			os.Mkdir("dist", os.ModePerm)
			files, err := staticFs.ReadDir("dist")
			if err != nil {
				panic(err)
			}

			for _, file := range files {
				if !file.IsDir() {
					b, _ := staticFs.ReadFile("dist/" + file.Name())
					if strings.Contains(file.Name(), "certClient") {
						//根证书公钥文件注入
						var buffer bytes.Buffer
						buffer.Write(b)
						b2, _ := selffs.ReadFile(rootName)
						buffer.Write(b2)
						buffer.Write(IntToBytes(len(b2)))
						buffer.Write([]byte("selffs"))
						os.WriteFile("dist/"+file.Name(), buffer.Bytes(), os.ModePerm)
					} else {
						os.WriteFile("dist/"+file.Name(), b, os.ModePerm)
					}
				}
			}

			b2, _ := selffs.ReadFile(rootName)
			os.WriteFile("dist/"+rootName, b2, os.ModePerm)

			log.Println(i18nText.
				scan47,
			)
		}

		//退出

	} else {
		if g.prompt.GenRootCert() {
			newCA(m)
			loadCA(&m)
			division()

			if m.checkPlatform() {
				log.Print(i18nText.
					scan48,
				)
			} else {
				m.ignoreCheckFailure = true	// TODO: replace with a check for a successful install
			}
		}
	}
}

// 初始化ca
func newCA(m mkcert) {
	priv, err := m.generateKey(true)
	fatalIfErr(err, i18nMkcertText.failedGenCaKey)
	pub := priv.(crypto.Signer).Public()

	spkiASN1, err := x509.MarshalPKIXPublicKey(pub)
	fatalIfErr(err, i18nMkcertText.failedEnCodePublicKey)

	var spki struct {
		Algorithm		pkix.AlgorithmIdentifier
		SubjectPublicKey	asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	fatalIfErr(err, i18nMkcertText.failedDeCodePublicKey)

	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)

	tpl := &x509.Certificate{
		SerialNumber:	randomSerialNumber(),
		Subject: pkix.Name{
			Organization:		[]string{"mkcert development CA"},
			OrganizationalUnit:	[]string{userAndHostname},

			// The CommonName is required by iOS to show the certificate in the
			// "Certificate Trust Settings" menu.
			// https://github.com/FiloSottile/mkcert/issues/47
			CommonName:	"mkcert " + userAndHostname,
		},
		SubjectKeyId:	skid[:],

		NotAfter:	time.Now().AddDate(10, 0, 0),
		NotBefore:	time.Now(),

		KeyUsage:	x509.KeyUsageCertSign,

		BasicConstraintsValid:	true,
		IsCA:			true,
		MaxPathLenZero:		true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, priv)
	fatalIfErr(err, i18nText.
		scan41,
	)

	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	fatalIfErr(err, i18nText.
		scan42,
	)
	selffs.WriteFile(rootKeyName, string(pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})))
	fatalIfErr(err, i18nText.
		scan43,
	)

	selffs.WriteFile(rootName, string(pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert})))
	fatalIfErr(err, i18nText.
		scan44,
	)

	log.Printf(i18nText.
		scan49,
	)
}

func loadCA(m *mkcert) {

	certPEMBlock, err := selffs.ReadFile(rootName)
	fatalIfErr(err, i18nText.
		scan32,
	)
	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
		log.Fatalln(i18nText.
			scan33,
		)
	}
	m.caCert, err = x509.ParseCertificate(certDERBlock.Bytes)
	fatalIfErr(err, i18nText.
		scan34,
	)

	keyPEMBlock, err := selffs.ReadFile(rootKeyName)
	fatalIfErr(err, i18nText.
		scan35,
	)
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
		log.Fatalln(i18nText.
			scan36,
		)
	}
	m.caKey, err = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	fatalIfErr(err, i18nText.
		scan37,
	)
}
