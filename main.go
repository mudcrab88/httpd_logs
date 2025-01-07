package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

var conf Config
var confFile = "config.json"

func findStr(pattern string, line string) string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return re.FindString(line)
}

func parseLine(line string, conf *Config) string {
	ip := findStr(conf.RegularExpressionIP, line)
	timeStr := findStr(conf.RegularExpressionTime, line)
	time, err := time.Parse(conf.TimeParseFormat, timeStr)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%s-%s", time.Format(conf.TimeReturnFormat), ip)
}

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

func parseFile(resultMap *map[string]int, file *os.File, mutex *sync.Mutex, conf *Config) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		mutex.Lock()
		addRecord(resultMap, parseLine(line, conf))
		mutex.Unlock()
	}
}

func parseLogsDir(path string, conf *Config) {
	os.Chdir(path)
	var resultMap = map[string]int{}
	var mutex sync.Mutex
	var wg sync.WaitGroup

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		file, err := os.Open(entry.Name())
		defer file.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		wg.Add(1)
		fmt.Println("Обработка файла: " + entry.Name())
		go func(file *os.File) {
			defer wg.Done()
			parseFile(&resultMap, file, &mutex, conf)
			fmt.Println("Закончена обработка файла: " + entry.Name())
		}(file)
	}
	wg.Wait()

	fmt.Println(len(resultMap))
	getMax(resultMap)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Не указаны аргументы командной строки")
		return
	}

	conf, err := NewConfig(confFile)

	if err != nil {
		fmt.Println("Ошибка при чтении конфигурации:", err)
	}

	parseLogsDir(os.Args[1], conf)
}
