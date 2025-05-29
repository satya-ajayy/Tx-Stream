package utils

import (
	// Go Internal Packages
	"strconv"
	"strings"
)

func JoinInt32Slice(ints []int32) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = strconv.FormatInt(int64(v), 10)
	}
	return strings.Join(strs, ",")
}
