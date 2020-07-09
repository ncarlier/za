package helper

import "regexp"

var botRegexp = regexp.MustCompile("(?i)bot|crawl|spider|robot|crawling")

// IsBotUserAgent test if User-Agent string is a bot
func IsBotUserAgent(name string) bool {
	return botRegexp.MatchString(name)
}
