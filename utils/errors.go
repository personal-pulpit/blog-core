package utils

import "strings"

func CheckErrorForWord(err error, word string) bool {
    return err != nil && strings.Contains(err.Error(), word)
}
