package helpers

import "strconv"

func StringToInt(str string)int{
	i,err := strconv.Atoi(str)
	if err != nil {
		//TODO:manage panics
		panic(err)
	}

	return i
}