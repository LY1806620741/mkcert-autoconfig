package main

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18nText struct {
	help, scan10, scan11, scan12, scan13, scan14, scan15, scan16, scan17, scan18, scan19, scan20, scan21, scan22, scan23, scan24, scan25, scan26, scan27, scan28, scan29, scan30, scan31, scan32, scan33, scan34, scan35, scan36, scan37, scan38, scan39, scan40, scan41, scan42, scan43, scan44, scan45, scan46, scan47, scan48, scan49, scan50, scan51, scan52, scan53, scan54, scan55, scan56, scan57, scan58, scan59, scan60, scan61, scan62, scan63, scan64, scan65, scan66, scan67, scan68, scan69, scan7, scan70, scan71, scan72, scan73, scan74, scan75, scan76, scan77, scan78, scan79, scan8, scan80, scan81, scan82, scan83, scan84, scan85, scan86, scan87, scan88, scan89, scan9, scan90, scan91, scan92, scan93, scan94, scan95 string
}

type I18nMkcertText struct {
	//
	moreOptions, shortUsage, advancedUsage, failedGenCaKey, failedEnCodePublicKey, failedDeCodePublicKey string
}

var localizer *i18n.Localizer

//go:embed active.*.toml
var LocaleFS embed.FS

var i18nText I18nText
var i18nMkcertText I18nMkcertText

func init() {
	userLanguage, _ := jibber_jabber.DetectLanguage()
	userTag := language.MustParse(userLanguage)
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "active.zh.toml")
	tag, _, _ := language.NewMatcher([]language.Tag{
		language.English,
		language.Chinese,
	}).Match(userTag)
	localizer = i18n.NewLocalizer(bundle, tag.String())

	i18nText.help = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "help",
			Other: `Usage of {{.Name}}:
		
    $ {{.Name}} auto
    use guide

    $ {{.Name}} mkcert
    use mkcert, the same as mkcert's usage.

`,
		}, TemplateData: map[string]interface{}{
			"Name": execNameWithOutSuffix,
		},
	})

	//初始化mkcert的说明文本
	i18nMkcertText.moreOptions = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "moreOptions",
			Other: `For more options, run "mkcert -help".`,
		},
	})

	i18nMkcertText.shortUsage = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "shortUsage",
			Other: `Usage of mkcert:

	$ mkcert -install
	Install the local CA in the system trust store.

	$ mkcert example.org
	Generate "example.org.pem" and "example.org-key.pem".

	$ mkcert example.com myapp.dev localhost 127.0.0.1 ::1
	Generate "example.com+4.pem" and "example.com+4-key.pem".

	$ mkcert "*.example.it"
	Generate "_wildcard.example.it.pem" and "_wildcard.example.it-key.pem".

	$ mkcert -uninstall
	Uninstall the local CA (but do not delete it).

`,
		},
	})
	i18nMkcertText.advancedUsage = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "advancedUsage",
			Description: "this is advancedUsage",
			Other: `Advanced options:

	-cert-file FILE, -key-file FILE, -p12-file FILE
		Customize the output paths.

	-client
		Generate a certificate for client authentication.

	-ecdsa
		Generate a certificate with an ECDSA key.

	-pkcs12
		Generate a ".p12" PKCS #12 file, also know as a ".pfx" file,
		containing certificate and key for legacy applications.

	-csr CSR
		Generate a certificate based on the supplied CSR. Conflicts with
		all other flags and arguments except -install and -cert-file.

	-CAROOT
		Print the CA certificate and key storage location.

	$CAROOT (environment variable)
		Set the CA certificate and key storage location. (This allows
		maintaining multiple local CAs in parallel.)

	$TRUST_STORES (environment variable)
		A comma-separated list of trust stores to install the local
		root CA into. Options are: "system", "java" and "nss" (includes
		Firefox). Autodetected by default.

`,
		},
	})

	i18nMkcertText.failedGenCaKey = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "failedGenCaKey",
			Other: `failed to generate the CA key`,
		},
	})
	i18nMkcertText.failedEnCodePublicKey = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "failedEnCodePublicKey",
			Other: `failed to encode public key`,
		},
	})
	i18nMkcertText.failedDeCodePublicKey = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "failedDeCodePublicKey",
			Other: `failed to decode public key`,
		},
	})

	//初始化自动收集的待翻译文本

	i18nText.scan30 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan30",
			Other: "failed to parse generated certificate",
		},
	})

	i18nText.scan41 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan41",
			Other: "failed to generate CA certificate",
		},
	})

	i18nText.scan44 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan44",
			Other: "failed to save CA certificate",
		},
	})

	i18nText.scan47 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan47",
			Other: "已生成授信客户端，在当前dist目录下，请在你的服务器进行部署",
		},
	})

	i18nText.scan53 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan53",
			Other: "ERROR: can't specify extra arguments when using -csr",
		},
	})

	i18nText.scan55 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan55",
			Other: "failed to create the CAROOT",
		},
	})

	i18nText.scan95 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan95",
			Other: "delete cert",
		},
	})

	i18nText.scan8 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan8",
			Other: "failed to generate certificate key",
		},
	})

	i18nText.scan23 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan23",
			Other: "   Warning: many browsers don't support second-level wildcards like %q ⚠️",
		},
	})

	i18nText.scan37 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan37",
			Other: "failed to parse the CA key",
		},
	})

	i18nText.scan57 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan57",
			Other: "Note: the local CA is not installed in the %s trust store.",
		},
	})

	i18nText.scan66 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan66",
			Other: `Warning: "certutil" is not available, so the CA can't be automatically installed in %s! ⚠️`,
		},
	})

	i18nText.scan82 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan82",
			Other: "failed to parse trust settings",
		},
	})

	i18nText.scan12 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan12",
			Other: "failed to save certificate",
		},
	})

	i18nText.scan21 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan21",
			Other: "\nCreated a new certificate valid for the following names 📜",
		},
	})

	i18nText.scan50 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan50",
			Other: "ERROR: you can't set -[un]install and -CAROOT at the same time",
		},
	})

	i18nText.scan54 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan54",
			Other: "ERROR: failed to find the default CA location, set one as the CAROOT env var",
		},
	})

	i18nText.scan70 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan70",
			Other: `Warning: "keytool" is not available, so the CA can't be automatically installed in Java's trust store! ⚠️`,
		},
	})

	i18nText.scan84 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan84",
			Other: "failed to serialize trust settings",
		},
	})

	i18nText.scan91 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan91",
			Other: "Note that if you never started %s, you need to do that at least once.",
		},
	})

	i18nText.scan94 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan94",
			Other: "add cert",
		},
	})

	i18nText.scan58 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan58",
			Other: "Note: the local CA is not installed in the Java trust store.",
		},
	})

	i18nText.scan80 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan80",
			Other: "failed to create temp file",
		},
	})

	i18nText.scan85 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan85",
			Other: "failed to write trust settings",
		},
	})

	i18nText.scan49 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan49",
			Other: "生成了根证书客户端automkcert-root 💥\n",
		},
	})

	i18nText.scan29 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan29",
			Other: "invalid CSR signature",
		},
	})

	i18nText.scan31 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan31",
			Other: "\nThe certificate is at \"%s\" ✅\n\n",
		},
	})

	i18nText.scan52 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan52",
			Other: "ERROR: can only combine -csr with -install and -cert-file",
		},
	})

	i18nText.scan62 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan62",
			Other: "The local CA is now installed in the system trust store! ⚡️",
		},
	})

	i18nText.scan81 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan81",
			Other: "failed to read trust settings",
		},
	})

	i18nText.scan86 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan86",
			Other: "Installing to the system store is not yet supported on this Linux 😣 but %s will still work.",
		},
	})

	i18nText.scan92 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan92",
			Other: "decode pem",
		},
	})

	i18nText.scan9 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan9",
			Other: "failed to generate certificate",
		},
	})

	i18nText.scan32 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan32",
			Other: "failed to read the CA certificate",
		},
	})

	i18nText.scan46 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan46",
			Other: "已导出到当前目录下",
		},
	})

	i18nText.scan56 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan56",
			Other: "Note: the local CA is not installed in the system trust store.",
		},
	})

	i18nText.scan65 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan65",
			Other: `Note: %s support is not available on your platform. ℹ️`,
		},
	})

	i18nText.scan74 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan74",
			Other: `Warning: "keytool" is not available, so the CA can't be automatically uninstalled from Java's trust store (if it was ever installed)! ⚠️`,
		},
	})

	i18nText.scan22 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan22",
			Other: " - %q",
		},
	})

	i18nText.scan26 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan26",
			Other: "failed to read the CSR",
		},
	})

	i18nText.scan39 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan39",
			Other: "failed to encode public key",
		},
	})

	i18nText.scan40 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan40",
			Other: "failed to decode public key",
		},
	})

	i18nText.scan43 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan43",
			Other: "failed to save CA key",
		},
	})

	i18nText.scan71 = ""

	i18nText.scan93 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan93",
			Other: "open root store",
		},
	})

	i18nText.scan25 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan25",
			Other: "failed to generate serial number",
		},
	})

	i18nText.scan35 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan35",
			Other: "failed to read the CA key",
		},
	})

	i18nText.scan16 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan16",
			Other: "\nThe certificate and key are at \"%s\" ✅\n\n",
		},
	})

	i18nText.scan15 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan15",
			Other: "failed to save PKCS#12",
		},
	})

	i18nText.scan17 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan17",
			Other: "\nThe certificate is at \"%s\" and the key at \"%s\" ✅\n\n",
		},
	})

	i18nText.scan45 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan45",
			Other: "Created a new local CA 💥\n",
		},
	})

	i18nText.scan83 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan83",
			Other: "ERROR: unsupported trust settings version:",
		},
	})

	i18nText.scan90 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan90",
			Other: "Installing in %s failed. Please report the issue with details about your environment at https://github.com/FiloSottile/mkcert/issues/new 👎",
		},
	})

	i18nText.scan20 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan20",
			Other: "It will expire on %s 🗓\n\n",
		},
	})

	i18nText.scan64 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan64",
			Other: "The local CA is now installed in the %s trust store (requires browser restart)! 🦊",
		},
	})

	i18nText.scan73 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan73",
			Other: `You can install "certutil" with "%s" and re-run "mkcert -uninstall" 👈`,
		},
	})

	i18nText.scan19 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan19",
			Other: "\nThe legacy PKCS#12 encryption password is the often hardcoded default \"changeit\" ℹ️\n\n",
		},
	})

	i18nText.scan38 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan38",
			Other: "failed to generate the CA key",
		},
	})

	i18nText.scan48 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan48",
			Other: "The local CA is already installed in the system trust store! 👍",
		},
	})

	i18nText.scan60 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan60",
			Other: "ERROR: %q is not a valid hostname, IP, URL or email: %s",
		},
	})

	i18nText.scan63 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan63",
			Other: "The local CA is already installed in the %s trust store! 👍",
		},
	})

	i18nText.scan75 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan75",
			Other: "The local CA is now uninstalled from the system trust store(s)! 👋",
		},
	})

	i18nText.scan76 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan76",
			Other: "The local CA is now uninstalled from the %s trust store(s)! 👋",
		},
	})

	i18nText.scan77 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan77",
			Other: "ERROR: %s: %s",
		},
	})

	i18nText.scan27 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan27",
			Other: "ERROR: failed to read the CSR: unexpected content",
		},
	})

	i18nText.scan51 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan51",
			Other: "ERROR: you can't set -install and -uninstall at the same time",
		},
	})

	i18nText.scan79 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan79",
			Other: `Warning: "sudo" is not available, and mkcert is not running as root. The (un)install operation might fail. ⚠️`,
		},
	})

	i18nText.scan89 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan89",
			Other: "ERROR: no %s security databases found",
		},
	})

	i18nText.scan28 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan28",
			Other: "failed to parse the CSR",
		},
	})

	i18nText.scan13 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan13",
			Other: "failed to save certificate key",
		},
	})

	i18nText.scan69 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan69",
			Other: "The local CA is now installed in Java's trust store! ☕️",
		},
	})

	i18nText.scan87 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan87",
			Other: "You can also manually install the root certificate at %q.",
		},
	})

	i18nText.scan10 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan10",
			Other: "failed to encode certificate key",
		},
	})

	i18nText.scan11 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan11",
			Other: "failed to save certificate and key",
		},
	})

	i18nText.scan14 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan14",
			Other: "failed to generate PKCS#12",
		},
	})

	i18nText.scan36 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan36",
			Other: "ERROR: failed to read the CA key: unexpected content",
		},
	})

	i18nText.scan7 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan7",
			Other: "ERROR: can't create new certificates because the CA key (rootCA-key.pem) is missing",
		},
	})

	i18nText.scan33 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan33",
			Other: "ERROR: failed to read the CA certificate: unexpected content",
		},
	})

	i18nText.scan42 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan42",
			Other: "failed to encode CA key",
		},
	})

	i18nText.scan67 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan67",
			Other: `Install "certutil" with "%s" and re-run "mkcert -install" 👈`,
		},
	})

	i18nText.scan72 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan72",
			Other: `Warning: "certutil" is not available, so the CA can't be automatically uninstalled from %s (if it was ever installed)! ⚠️`,
		},
	})

	i18nText.scan78 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan78",
			Other: "ERROR: failed to execute \"%s\": %s\n\n%s\n",
		},
	})

	i18nText.scan88 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan88",
			Other: "failed to read root certificate",
		},
	})

	i18nText.scan18 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan18",
			Other: "\nThe PKCS#12 bundle is at \"%s\" ✅\n",
		},
	})

	i18nText.scan34 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan34",
			Other: "failed to parse the CA certificate",
		},
	})

	i18nText.scan59 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan59",
			Other: "Run \"mkcert -install\" for certificates to be trusted automatically ⚠️",
		},
	})

	i18nText.scan61 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan61",
			Other: "ERROR: %q is not a valid hostname, IP, URL or email",
		},
	})

	i18nText.scan68 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan68",
			Other: "The local CA is already installed in Java's trust store! 👍",
		},
	})

	i18nText.scan24 = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "scan24",
			Other: "\nReminder: X.509 wildcards only go one level deep, so this won't match a.b.%s ℹ️",
		},
	})

}

func (i *I18nText) errUnknownGroupCommand(groupCommand string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "errUnknownGroupCommand",
			Other: `err. Unknown command "{{.Name}}"`,
		}, TemplateData: map[string]interface{}{
			"Name": execNameWithOutSuffix + " " + groupCommand,
		},
	})
}
