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

type Parser struct {
}

func newParser() (*Parser, error) {
	var parser Parser

	return &parser, nil
}

func (parser *Parser) findStr(pattern string, line string) string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return re.FindString(line)
}

func (parser *Parser) parseLine(line string, conf *Config) string {
	ip := parser.findStr(conf.RegularExpressionIP, line)
	timeStr := parser.findStr(conf.RegularExpressionTime, line)
	time, err := time.Parse(conf.TimeParseFormat, timeStr)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%s-%s", time.Format(conf.TimeReturnFormat), ip)
}

func (parser *Parser) parseFile(resultMap *map[string]int, file *os.File, mutex *sync.Mutex, conf *Config) {
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		mutex.Lock()
		addRecord(resultMap, parser.parseLine(line, conf))
		mutex.Unlock()
	}
}

func (parser *Parser) parseLogsDir(path string, conf *Config) {
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
			parser.parseFile(&resultMap, file, &mutex, conf)
			fmt.Println("Закончена обработка файла: " + entry.Name())
		}(file)
	}
	wg.Wait()

	fmt.Println(len(resultMap))
	getMax(resultMap)
}
