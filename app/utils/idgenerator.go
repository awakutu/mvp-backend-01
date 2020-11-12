package utils

import "math/rand"

func RangeIn(low, hi int) int {
	//kita harus perbarui seed dari fungsi random ini
	//kalau tidak, setiap server dimulai, seed nya akan sama, mengakibatkan angka random pertama yg dihasilkan selalu sama

	return low + rand.Intn(hi-low)
}
