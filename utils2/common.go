package utils2

import (
	"fmt"
	"strconv"
)

func In(obj interface{}, objList []interface{}) bool {
	for _, v := range objList {
		if obj == v {
			return true
		}
	}
	return false
}

func ENum2float64(enum interface{}) float64 {
	var newNum float64
	value := enum.(string)

	_, err := fmt.Sscanf(value, "%e", &newNum)
	if err != nil {
		fmt.Printf("%v not convert to int", enum)
	}
	return newNum
}

func String2int(str interface{}) int {
	strNumber, ok := str.(string)
	if !ok {
		fmt.Printf("%v type is not string", str)
	}
	number, err := strconv.Atoi(strNumber)
	if err == nil {
		return number
	}
	return 0
}
