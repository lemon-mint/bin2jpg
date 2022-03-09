package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/lemon-mint/bin2jpg"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

func printUsage() {
	fmt.Println("Usage: jpg2bin [options] inputfile")
	fmt.Println("\t-o: output file name")
	os.Exit(1)
}

func main() {
	var flags = map[string]string{}
	args := os.Args[1:]
loop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			flags["o"] = args[i+1]
			i++
		default:
			if i == len(args)-1 {
				args = args[i:]
				break loop
			}
			fmt.Println("Unknown Argument:", args[i])
			printUsage()
		}
	}
	var inputFileName, outputFileName string
	if len(args) > 0 {
		inputFileName = args[0]
		inputFileName = strings.Trim(inputFileName, " ")
		if o, ok := flags["o"]; ok {
			outputFileName = o
		} else {
			switch {
			case strings.HasSuffix(inputFileName, ".jpg"):
				outputFileName = inputFileName[:len(inputFileName)-4]
			case strings.HasSuffix(inputFileName, ".jpeg"):
				outputFileName = inputFileName[:len(inputFileName)-5]
			case strings.HasSuffix(inputFileName, ".png"):
				outputFileName = inputFileName[:len(inputFileName)-4]
			case strings.HasSuffix(inputFileName, ".gif"):
				outputFileName = inputFileName[:len(inputFileName)-4]
			case strings.HasSuffix(inputFileName, ".webp"):
				outputFileName = inputFileName[:len(inputFileName)-5]
			default:
				outputFileName = inputFileName + ".bin"
			}
		}
	} else {
		printUsage()
	}

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer inputFile.Close()
	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	data, err := bin2jpg.ImageDecode(img)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
