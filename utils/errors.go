package utils

import (
	"fmt"
	"strings"
)

func CheckErrorForWord(err error, word string) bool {
    return err != nil && strings.Contains(err.Error(), word)
}
func GetValidationError(err error)string{
    return fmt.Sprintf("validation error:%s",err.Error())
}