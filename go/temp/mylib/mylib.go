package mylib

import "fmt"

func NonEmptyStrings(strs []string) []string {
	i := 0
	for _, str := range strs {
		if str == "" {
			continue
		}
		strs[i] = str
		i++
	}
	return strs[:i]
}

func Reverse(arr *[3]int) {
	fmt.Println("arr len: ", len(*arr))
	kvMap := make(map[string]int)
	kvMap["a"] = 1
	kvMap["b"] = 2

	kvMap1 := map[string]int{
		"hello": 1,
	}
	fmt.Println(len(kvMap1))
}
