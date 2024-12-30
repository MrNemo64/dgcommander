package util

import "github.com/bwmarrin/discordgo"

type Localizable[B any] struct {
	Upper         B
	Value         string
	Localizations map[discordgo.Locale]string
}

func (l *Localizable[B]) Set(value string) B {
	l.Value = value
	return l.Upper
}

func (l *Localizable[B]) SetLocalizations(localizations map[discordgo.Locale]string) B {
	l.Localizations = localizations
	return l.Upper
}
