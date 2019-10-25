package util

import (
	"strconv"
)

func Eth2float64(str []string) ([]*string, error) {

	if len(str) > 0 {
		var str11 []*string
		for i := 0; i < len(str); i++ {

			f, _ := strconv.ParseFloat(str[i], 64)

			a := f / 1000000000000000000.00 //float
			c := strconv.FormatFloat(a, 'f', 10, 64)
			str11 = append(str11, &c)

		}
		return str11, nil

	}
	return nil, nil

}
