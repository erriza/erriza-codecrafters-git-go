package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":		
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}
		
		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}
		
		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "usage: got cat-file -p <object-hash>\n")
			os.Exit(1)
		}

		readFlag := os.Args[2]
		objectHash := os.Args[3]

		if readFlag != "-p" && len(objectHash) != 40 {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object-hash>\n")
			os.Exit(1)
		}

		dirName := objectHash[:2]
		fileName := objectHash[2:]

		filePath := fmt.Sprintf("./.git/objects/%s/%s", dirName, fileName)

		fileContents, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		}

		b := bytes.NewReader(fileContents)
		r, err := zlib.NewReader(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decompressing the file: %s\n", err)
			os.Exit(1)
		}

		decompressedData, err := io.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading decompressed data: %s\n", err)
			os.Exit(1)
		}
		r.Close()

		//Find the idx of the null terminaron
		nullIndex := bytes.IndexByte(decompressedData, 0)
		if nullIndex == -1 {
			fmt.Fprintf(os.Stderr, "Invalid object format: missing metadata separator\n")
			os.Exit(1)
		}

		content := decompressedData[nullIndex+1:]
		fmt.Print(string(content))

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
