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

/*	if len(str) == 1 {
		feet, err := strconv.ParseFloat(*str[0], 64)
		if err != nil {
			return str, fmt.Errorf("转换出错：%v\n", err.Error())
		} else {

			a := feet / 1000000000000000000.00

			return a, nil

		}

	} else {

		for key, _ := range str {
			var str1 []*string
			feet, err := strconv.ParseFloat(str[key], 64)

			if err != nil {
				return str, fmt.Errorf("转换出错：%v\n", err.Error())
			} else {

				a := feet / 1000000000000000000.00
				float := strconv.FormatFloat(a, 'f', 10, 64)
				str1 := append(str1, &float)
				return str1, nil

			}
		}
	}*/
