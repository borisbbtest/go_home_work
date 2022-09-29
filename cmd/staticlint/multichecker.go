package main

import (
	mainmultichecker "github.com/borisbbtest/go_home_work/internal/multichecker"
)

func main() {
	mc := mainmultichecker.InitMultichecker()
	mc.Start()
}
