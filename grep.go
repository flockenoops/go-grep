package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var patternStr = flag.String("p", "", "Pattern to match - mandatory.")
var fileStr = flag.String("f", "", "File path, mutually exclusive with -t. If both are empty, pipe is used.")
var textStr = flag.String("t", "", "Text string, mutually exclusive with -f. If both are empty, pipe is used.")
var regexpBool = flag.Bool("r", false, "Enables regexp matching, pattern must be a valid regular expression.")
var colorBool = flag.Bool("c", false, "Colorize matches.")

type lineReadyToPrint struct {
	splitLine     []string
	matches       []string
	isMatchSuffix bool
}

func getScanner() bufio.Scanner {
	var scanner bufio.Scanner
	if *textStr == "" && *fileStr == "" {
		file, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		if file.Mode()&os.ModeNamedPipe == 0 {
			flag.Usage()
			log.Fatal("Couldn't read input.")
		} else {
			scanner = *bufio.NewScanner(os.Stdin)
		}
	} else if *fileStr != "" && *textStr != "" {
		flag.Usage()
		log.Fatal("-f and -t cannot be used in conjunction.")
	} else if *fileStr == "" {
		scanner = *bufio.NewScanner(strings.NewReader(*textStr))
	} else {
		file, err := os.Open(*fileStr)
		if err != nil {
			log.Fatal(err)
		}
		scanner = *bufio.NewScanner(file)
	}
	return scanner
}

func regexpMatch(pattern string, scanner bufio.Scanner) []lineReadyToPrint {
	var linesToReturn []lineReadyToPrint
	matcher := regexp.MustCompile(pattern)
	for scanner.Scan() {
		var lineResult lineReadyToPrint
		line := scanner.Text()
		if matcher.MatchString(line) {
			lineResult.splitLine = matcher.Split(line, -1)
			lineResult.matches = matcher.FindAllString(line, -1)
			if len(lineResult.matches) == len(lineResult.splitLine) {
				lineResult.isMatchSuffix = true
			}
			linesToReturn = append(linesToReturn, lineResult)
		}
	}
	return linesToReturn
}

func stringMatch(pattern string, scanner bufio.Scanner) []lineReadyToPrint {
	var linesToReturn []lineReadyToPrint
	var lineResult lineReadyToPrint
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			lineResult.splitLine = strings.Split(line, pattern)
			lineResult.isMatchSuffix = strings.HasSuffix(line, pattern)
			linesToReturn = append(linesToReturn, lineResult)
		}
	}
	return linesToReturn
}

func regexpPrint(lines []lineReadyToPrint) {
	for l := 0; l < len(lines); l++ {
		if *colorBool {
			for m := range lines[l].matches {
				lines[l].matches[m] = "\033[1;36m" + lines[l].matches[m] + "\033[0m"
			}
		}
		var lineToPrint string
		for s := 0; s < len(lines[l].splitLine)-1; s++ {
			lineToPrint = lineToPrint + lines[l].splitLine[s] + lines[l].matches[s]
		}
		lineToPrint = lineToPrint + lines[l].splitLine[len(lines[l].splitLine)-1]
		fmt.Println(lineToPrint)
	}
}

func stringPrint(lines []lineReadyToPrint) {
	var pattern string
	if *colorBool {
		pattern = "\033[1;36m" + *patternStr + "\033[0m"
	} else {
		pattern = *patternStr
	}
	for l := 0; l < len(lines); l++ {
		var lineToPrint string
		for s := 0; s < len(lines[l].splitLine)-1; s++ {
			lineToPrint = lineToPrint + lines[l].splitLine[s] + pattern
		}
		if !lines[l].isMatchSuffix {
			lineToPrint = lineToPrint + lines[l].splitLine[len(lines[l].splitLine)-1]
		}
		fmt.Println(lineToPrint)
	}
}

func match(pattern string, scanner bufio.Scanner) {
	var allMatchedLines []lineReadyToPrint
	if *regexpBool {
		allMatchedLines = regexpMatch(pattern, scanner)
		regexpPrint(allMatchedLines)
	} else {
		allMatchedLines = stringMatch(pattern, scanner)
		stringPrint(allMatchedLines)
	}
}

func main() {
	flag.Parse()

	if *patternStr == "" {
		flag.Usage()
		log.Fatal("Please provide a search pattern.")
	}

	scanner := getScanner()

	match(*patternStr, scanner)
}
