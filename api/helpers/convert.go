package helpers

import "strconv"

func StringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return i, nil
}
