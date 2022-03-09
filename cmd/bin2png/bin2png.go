package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/lemon-mint/bin2jpg"
)

func printUsage() {
	fmt.Println("Usage: bin2png [options] inputfile")
	fmt.Println("\t-o: output file name")
	fmt.Println("\t-e: output file format (png | jpg)")
	os.Exit(1)
}

func main() {
	var flags = map[string]string{}
	flags["e"] = "png"
	args := os.Args[1:]
loop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			flags["o"] = args[i+1]
			i++
		case "-e":
			flags["e"] = args[i+1]
			i++
			if flags["e"] != "png" && flags["e"] != "jpg" {
				fmt.Println("Error: -e must be png or jpg")
				printUsage()
			}
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
		if o, ok := flags["o"]; ok {
			outputFileName = o
		} else {
			outputFileName = inputFileName + "." + flags["e"]
		}
	} else {
		printUsage()
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	data, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	img := bin2jpg.ImageEncode(data)

	switch flags["e"] {
	case "png":
		err = png.Encode(outputFile, img)
	case "jpg":
		err = jpeg.Encode(outputFile, img, nil)
	}
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
