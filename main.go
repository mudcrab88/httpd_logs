package main

import (
	"fmt"
	"os"
)

var conf Config
var confFile = "config.json"

func addRecord(res *map[string]int, str string) {
	if val, ok := (*res)[str]; ok {
		(*res)[str] = val + 1
	} else {
		(*res)[str] = 1
	}
}

func getMax(resultMap map[string]int) {
	maxKey := ""
	maxValue := 0
	for key, value := range resultMap {
		if value > maxValue {
			maxKey = key
			maxValue = value
			fmt.Println(maxKey, maxValue)
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Не указаны аргументы командной строки")
		return
	}

	conf, err := newConfig(confFile)

	if err != nil {
		fmt.Println("Ошибка при чтении конфигурации:", err)
	}

	parser, err := newParser()

	parser.parseLogsDir(os.Args[1], conf)
}
