package utils

func MakeTitle(content string) string {

	if len(content) <= 30 {
		return content
	}

	return content[:30] + "..."
}
