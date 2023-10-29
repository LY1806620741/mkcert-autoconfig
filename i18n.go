package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var userPrefs = []language.Tag{}

var serverLangs = []language.Tag{
	language.SimplifiedChinese, // zh-Hans
	language.AmericanEnglish,   // en-US
}

type i18n struct {
}

var matcher = language.NewMatcher(serverLangs)

func init() {
}

func (i *i18n) currentLanguage() string {
	tag, _, _ := matcher.Match(userPrefs...)
	return display.English.Tags().Name(tag)
}
