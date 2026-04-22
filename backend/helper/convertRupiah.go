package helper

import "fmt"

func FormatRupiah(amount float64) string {
	intAmount := int64(amount)
	str := fmt.Sprintf("%d", intAmount)
	result := ""
	for i, c := range reverse(str) {
		if i > 0 && i%3 == 0 {
			result = "." + result
		}
		result = string(c) + result
	}
	return "Rp " + result
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
