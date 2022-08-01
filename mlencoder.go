package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	R0 int    `yaml:"r0"`
	P1 int    `yaml:"p1"`
	P2 int    `yaml:"p2"`
	CT string `yaml:"ct"`
}

func encodeLine(line string, r *int, p1 int, p2 int, ct string) (string, error) {
	var sb strings.Builder
	for _, ch := range line {
		var k int
		if ch == '.' {
			k = 0
		} else if ch == '-' {
			k = 1
		} else if ch >= 'a' && ch <= 'z' {
			k = int(ch-'a') + 2
		} else if ch >= 'A' && ch <= 'Z' {
			k = int(ch-'A') + 28
		} else if ch >= '0' && ch <= '9' {
			k = int(ch-'0') + 54
		} else {
			k = -1
		}
		if k >= 0 {
			k = (*r & 0x3f) ^ k
			sb.WriteByte(ct[k])
		} else {
			sb.WriteRune(ch)
		}
		*r = p1*(*r) + p2
	}
	*r = p1*(*r) + p2
	*r = p1*(*r) + p2
	return sb.String(), nil
}

func encodeFile(inputFilePath string, outputFilePath string, config *Config) error {
	r := config.R0
	readFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer readFile.Close()
	var writeFile *os.File
	writeFile, err = os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer writeFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) > 0 {
			encodedLine, err := encodeLine(line, &r, config.P1, config.P2, config.CT)
			if err != nil {
				return err
			}
			if _, err = writeFile.WriteString(encodedLine + "\r\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	var config Config
	configFilePath := "mlencoder.yml"
	var inputFilePath string
	var outputFilePath string
	if len(os.Args) == 5 && os.Args[1] == "-config" {
		configFilePath = os.Args[2]
		inputFilePath = os.Args[3]
		outputFilePath = os.Args[4]
	} else if len(os.Args) == 3 {
		inputFilePath = os.Args[1]
		outputFilePath = os.Args[2]
	} else {
		fmt.Println("USAGE: mlencoder [-config config_file_path] input_file_path output_file_path")
		return
	}
	configFile, err := ioutil.ReadFile(configFilePath)
	if err == nil {
		err = yaml.Unmarshal(configFile, &config)
	}
	if err != nil {
		log.Fatalf("Read config: %v", err)
	}
	err = encodeFile(inputFilePath, outputFilePath, &config)
	if err != nil {
		log.Fatalf("Encode file: %v", err)
	}
}
