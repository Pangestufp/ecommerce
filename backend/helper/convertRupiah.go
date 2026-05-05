package helper

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func FormatRupiah(amount decimal.Decimal) string {
	// pisah integer dan desimal
	intAmount := amount.IntPart()
	str := fmt.Sprintf("%d", intAmount)
	result := ""
	for i, c := range reverse(str) {
		if i > 0 && i%3 == 0 {
			result = "." + result
		}
		result = string(c) + result
	}

	// tambah koma + desimal kalau ada
	cents := amount.Sub(decimal.NewFromInt(intAmount))
	if !cents.IsZero() {
		result += "," + cents.StringFixed(2)[2:] // ambil 2 digit setelah titik
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
