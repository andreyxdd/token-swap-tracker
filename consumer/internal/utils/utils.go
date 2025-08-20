package utils

import "fmt"

func BuildTokenPairkey(t1, t2 string) string {
	return fmt.Sprintf("%s-%s", t1, t2)
}

func Contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
