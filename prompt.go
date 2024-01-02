package main

//go:generate mockgen -destination mock_prompt_test.go -package main -source prompt.go Prompt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

type Prompt interface {
	GenRootCert() bool
	RootMenu() int
	InputHost() []string
}

type prompt struct {
}

type item struct {
	Name		string
	Description	string
}

func (p *prompt) GenRootCert() bool {
	items := []item{
		{Name: "生成网站根证书(不覆盖文件)", Description: "检查到" + execNameWithOutSuffix + "未生成根证书，根证书是颁发https证书所依赖的,你需要先生成它"},
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
		Label:		"你要做什么",
		Items:		items,
		Templates:	templates,
		Size:		4,
		Searcher:	searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}
	return i == 0
}

func (p *prompt) RootMenu() int {
	items := []item{
		{Name: "签名新证书", Description: "当前具备根证书，可以签名子证书"},
		{Name: "导出根证书", Description: "当前携带根证书公私钥，具有完全权限，可导出公私钥，请妥善保管"},
		{Name: "生成证书授信端", Description: "生成携带根证书公钥的授信客户端,可以部署到服务器让其他人安装根证书公钥授信"},
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
		Label:		"当前是携带根证书,你要做什么",
		Items:		items,
		Templates:	templates,
		Size:		4,
		Searcher:	searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}
	return i
}

func (p *prompt) InputHost() []string {
	validate := func(input string) error {
		hosts := strings.Split(input, ",")
		ipv4Regex, _ := regexp.Compile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
		hostRegex, _ := regexp.Compile(`^(?:[\w-]+\.)+[\w-]+$`)
		for idx, host := range hosts {
			match := ipv4Regex.MatchString(host)
			if !match {
				match = hostRegex.MatchString(host)
				if !match {
					return errors.New("第" + fmt.Sprint(idx+1) + "个域名或ip的格式不正确:" + host)
				}
			}
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:		"{{ . }} ",
		Valid:		"{{ . | green }} ",
		Invalid:	"{{ . | red }} ",
		Success:	"{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:		"输入服务器的域名以及ip,逗号隔开",
		Templates:	templates,
		Validate:	validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return []string{}
	}

	fmt.Printf("You answered %s\n", result)
	return []string{result}
}
