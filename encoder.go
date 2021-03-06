package main

import (
	"fmt"
	"strconv"
	"strings"
)

func encode(userId string, rid int32) string {
	return fmt.Sprintf("%s_%d", userId, rid)
}

func decode(key string) (userId string, rid int32, err error) {
	var (
		idx int
		t   int64
	)
	if idx = strings.IndexByte(key, '_'); idx == -1 {
		err = ErrDecodeKey
		return
	}
	userId = key[:idx];
	if t, err = strconv.ParseInt(key[idx+1:], 10, 32); err != nil {
		return
	}
	rid = int32(t)
	return
}
