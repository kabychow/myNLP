package model

import "strings"

type Tokens []string

func (tokens Tokens) String() string {
	return strings.Join(tokens, "|")
}
