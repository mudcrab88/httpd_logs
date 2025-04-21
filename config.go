package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	TimeReturnFormat      string `json:"time_return_format"`
	TimeParseFormat       string `json:"time_parse_format"`
	RegularExpressionIP   string `json:"regular_expression_ip"`
	RegularExpressionTime string `json:"regular_expression_time"`
}

func newConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling config:", err)
		return nil, err
	}

	return &config, nil
}
