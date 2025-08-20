package utils

import "fmt"

func BuildTokenPairkey(t1, t2 string) string {
	return fmt.Sprintf("%s-%s", t1, t2)
}
