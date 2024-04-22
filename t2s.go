package main

import (
 "flag"
 "fmt"
 "io/ioutil"
 "os"
 "path/filepath"
 "strings"
)

// readFile reads the contents of a file and returns it as a string.
func readFile(filename string) (string, error) {
 content, err := ioutil.ReadFile(filename)
 if err != nil {
  return "", err
 }
 return string(content), nil
}

// writeFile writes the contents to a file.
func writeFile(filename string, content string) error {
 err := ioutil.WriteFile(filename, []byte(content), 0644)
 if err != nil {
  return err
 }
 return nil
}

// convertSpacesToTabs converts 4 spaces to 1 tab in the given text.
func convertSpacesToTabs(text string) string {
 result := strings.ReplaceAll(text,"\t", "    ")
 return result
}

// convertIndentation converts indentation of a file from spaces to tabs.
func convertIndentation(filename string) error {
 content, err := readFile(filename)
 if err != nil {
  return err
 }

 convertedContent := convertSpacesToTabs(content)

 if err := writeFile(filename, convertedContent); err != nil {
  return err
 }

 fmt.Println("Converted indentation of file:", filename)
 return nil
}

func main() {
 // Define command-line flags
 dirPtr := flag.String("dir", ".", "Directory to process")
 flag.Parse()

 directory := *dirPtr

 // Walk through files in the directory
 err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
  if err != nil {
   return err
  }
  // Check if it's a regular file
  if !info.IsDir() {
   if err := convertIndentation(path); err != nil {
    fmt.Println("Error:", err)
   }
  }
  return nil
 })
 if err != nil {
  fmt.Println("Error:", err)
 }
}

