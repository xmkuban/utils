package utils

import "regexp"

var emojiRxg = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)

func FilterEmoji(source string) (result string) {
	if source == "" {
		return
	}
	return emojiRxg.ReplaceAllString(source, "")
}
