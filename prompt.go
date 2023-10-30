package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type prompt struct {
}

type item struct {
	Name        string
	Description string
}

func (p *prompt) GenRootCert() bool {
	items := []item{
		{Name: "生成网站根证书", Description: "检查到" + execNameWithOutSuffix + "未生成根证书，根证书是颁发https证书所依赖的,你需要先生成它"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
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
		Label:     "你要做什么",
		Items:     items,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}
	return i == 0
}
