package main

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18nText struct {
	//
	help string
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
	// println("Language:", userLanguage)
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
