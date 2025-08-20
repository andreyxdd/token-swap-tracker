package utils

import "fmt"

func BuildSemicolonKey(key, window string) string {
	return fmt.Sprintf("%s:%s", key, window)
}

func BuildHyphenKey(t1, t2 string) string {
	return fmt.Sprintf("%s-%s", t1, t2)
}

func Contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
