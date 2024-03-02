package main

import (
	"strconv"
	"time"
)

type Clock interface {
	Now() string;
}

type TimeClock struct {};

func (clock TimeClock) Now() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}