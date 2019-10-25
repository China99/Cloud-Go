package main

import "fmt"

func main() {
	//value := util.GetRealDecimalValue("4.111", 10)
	//eth := GetValETH([]string{"111111111111111111111"}, 2)
	value := []string{"1101111"}
	var decimal = 6
	if len(value) > 0 {
		for i := 0; i < len(value); i++ {

			fmt.Println(len(value[i]))
			strlen := len(value[i])
			s := value[i][:(strlen - decimal)]
			s2 := value[i][strlen-decimal : strlen]
			s3 := s + "." + s2
			fmt.Println(s3)
		}

	}
}

//小数点左移 decimal位
func GetValETH(value []string, decimal int) interface{} {
	//strlen := len(value)
	//if strlen>decimal{
	//	for i:=0;i<=strlen-decimal ;i++  {
	//		if i==strlen-decimal{
	//			fa
	//		}
	//	}
	//}
	return nil

}
