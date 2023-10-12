package util

import (
	"fmt"
	"strconv"
	"strings"
)

type QueryListFn func(string) (interface{}, error)

func QueryListToStrings(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("empty string")
	}
	return s, nil
}

func QueryListToInts64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

var QueryListToInts = strconv.Atoi

func Values[T any](q string, fn func(string) (T, error)) []T {
	res := make([]T, 0)

	s := strings.Split(q, ",")
	for _, v := range s {
		if i, err := fn(v); err == nil {
			res = append(res, i)
		}
	}

	if len(res) == 0 {
		return nil
	}
	return res
}
