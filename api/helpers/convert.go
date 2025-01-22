package helpers

import "strconv"

func StringToInt(str string)int{
	i,err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return i
}