package helper

import (
	"strings"
)

func TrimSpan(str string) string {
	str = strings.TrimSpace(str)
	if strings.Contains(str, "Spam\n") {
		str = strings.Replace(str, "Spam\n", "", -1)
	}
	if strings.Contains(str, "\n") {
		str = strings.Replace(str, "\n", "", -1)
	}
	if strings.Contains(str, " ") {
		str = strings.Replace(str, " ", "", -1)
	}
	if strings.Contains(str, ":") {
		str = strings.Replace(str, ":", "", -1)
	}

	if strings.Contains(str, "Canceled") {
		str = strings.Replace(str, "Canceled", "", -1)
	}

	if strings.Contains(str, "viaAPI.") {
		str = strings.Replace(str, "viaAPI.", "", -1)
	}

	if strings.Contains(str, "ParticipateDoNotInitiate") {
		str = strings.Split(str, "ParticipateDoNotInitiate")[1]
	}

	return str
}

func StringJoinString(val ...string) string {
	builder := strings.Builder{}
	for _, v := range val {
		builder.WriteString(v)
	}
	return builder.String()
}
