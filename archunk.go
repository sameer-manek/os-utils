package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run compressor.go <input_dir> <output_prefix> <max_zip_size>")
		return
	}

	inputDir := os.Args[1]
	outputPrefix := os.Args[2]
	maxZipSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid max ZIP size:", os.Args[3])
		return
	}

	err = compressDirectory(inputDir, outputPrefix, maxZipSize * 1024 * 1024)
	if err != nil {
		fmt.Println("Error compressing directory:", err)
		return
	}

	fmt.Println("Directory compressed successfully!")
}

func compressDirectory(inputDir, outputPrefix string, maxZipSize int) error {
	var currentZipSize int64
	var currentZipIndex int
	var currentZipFile *os.File
	var currentZipWriter *zip.Writer

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if currentZipWriter == nil || currentZipSize+int64(info.Size()) > int64(maxZipSize) {
			if currentZipWriter != nil {
				currentZipWriter.Close()
				currentZipFile.Close()
			}

			currentZipIndex++
			currentZipSize = 0
			currentZipFile, err = os.Create(fmt.Sprintf("%s_%03d.zip", outputPrefix, currentZipIndex))
			if err != nil {
				return err
			}
			currentZipWriter = zip.NewWriter(currentZipFile)
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(filepath.Base(inputDir), path[len(inputDir)+1:])

		writer, err := currentZipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		currentZipSize += int64(info.Size())
		return nil
	})

	if currentZipWriter != nil {
		currentZipWriter.Close()
		currentZipFile.Close()
	}

	return err
}


