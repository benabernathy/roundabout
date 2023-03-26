package internal

import _ "embed"

var (
	//go:embed resources/motd.ans
	defaultMotd string
)

func GetDefaultMotd() string {
	return defaultMotd
}
