package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"log"
	"time"
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
		loadCA(&m)
	} else {
		if g.prompt.GenRootCert() {
			newCA(m)
			loadCA(&m)
			division()

			if m.checkPlatform() {
				log.Print("The local CA is already installed in the system trust store! 👍")
			} else {
				m.ignoreCheckFailure = true // TODO: replace with a check for a successful install
			}
		}
	}
}

// 初始化ca
func newCA(m mkcert) {
	priv, err := m.generateKey(true)
	fatalIfErr(err, "failed to generate the CA key")
	pub := priv.(crypto.Signer).Public()

	spkiASN1, err := x509.MarshalPKIXPublicKey(pub)
	fatalIfErr(err, "failed to encode public key")

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	fatalIfErr(err, "failed to decode public key")

	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)

	tpl := &x509.Certificate{
		SerialNumber: randomSerialNumber(),
		Subject: pkix.Name{
			Organization:       []string{"mkcert development CA"},
			OrganizationalUnit: []string{userAndHostname},

			// The CommonName is required by iOS to show the certificate in the
			// "Certificate Trust Settings" menu.
			// https://github.com/FiloSottile/mkcert/issues/47
			CommonName: "mkcert " + userAndHostname,
		},
		SubjectKeyId: skid[:],

		NotAfter:  time.Now().AddDate(10, 0, 0),
		NotBefore: time.Now(),

		KeyUsage: x509.KeyUsageCertSign,

		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, priv)
	fatalIfErr(err, "failed to generate CA certificate")

	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	fatalIfErr(err, "failed to encode CA key")
	selffs.WriteFile(rootKeyName, string(pem.EncodeToMemory(
		&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})))
	fatalIfErr(err, "failed to save CA key")

	selffs.WriteFile(rootName, string(pem.EncodeToMemory(
		&pem.Block{Type: "CERTIFICATE", Bytes: cert})))
	fatalIfErr(err, "failed to save CA certificate")

	log.Printf("Created a new local CA 💥\n")
}

func loadCA(m *mkcert) {

	certPEMBlock, err := selffs.ReadFile(rootName)
	fatalIfErr(err, "failed to read the CA certificate")
	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
		log.Fatalln("ERROR: failed to read the CA certificate: unexpected content")
	}
	m.caCert, err = x509.ParseCertificate(certDERBlock.Bytes)
	fatalIfErr(err, "failed to parse the CA certificate")

	keyPEMBlock, err := selffs.ReadFile(rootKeyName)
	fatalIfErr(err, "failed to read the CA key")
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
		log.Fatalln("ERROR: failed to read the CA key: unexpected content")
	}
	m.caKey, err = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	fatalIfErr(err, "failed to parse the CA key")
}
