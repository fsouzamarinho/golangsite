package utils

import "strconv"

func AtoiOrZero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}