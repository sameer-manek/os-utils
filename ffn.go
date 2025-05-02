// simple function to parse and alter function names and their calls from snake case to camel case.
// it maintains the case of the first character tho

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var funcRenameMap = map[string]string{}

func main() {
	root := "./" // Change to desired path or pass as argument

	fmt.Println("🔍 First pass: Scanning for snake_case function declarations...")
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			collectRenameMap(path)
		}
		return nil
	})

	if len(funcRenameMap) == 0 {
		fmt.Println("✅ No snake_case function names found.")
		return
	}

	fmt.Println("\n🛠 Second pass: Rewriting function names across files...")
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			if err := updateFile(path); err != nil {
				fmt.Printf("❌ Error updating file %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func collectRenameMap(path string) {
	fmt.Printf("📄 Scanning: %s\n", path)
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, 0)
	if err != nil {
		fmt.Printf("  ⚠️ Skipping due to parse error: %v\n", err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name != nil {
			old := fn.Name.Name
			if strings.Contains(old, "_") {
				newName := toCamelCasePreserveFirst(old)
				if old != newName {
					funcRenameMap[old] = newName
					fmt.Printf("  🔁 Found: %s → %s\n", old, newName)
				}
			}
		}
		return true
	})
}

func updateFile(path string) error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	updated := false

	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncDecl:
			if newName, ok := funcRenameMap[v.Name.Name]; ok {
				fmt.Printf("  ✏️  Updating declaration: %s → %s in %s\n", v.Name.Name, newName, path)
				v.Name.Name = newName
				updated = true
			}
		case *ast.Ident:
			if newName, ok := funcRenameMap[v.Name]; ok {
				fmt.Printf("  ✏️  Updating usage: %s → %s in %s\n", v.Name, newName, path)
				v.Name = newName
				updated = true
			}
		case *ast.SelectorExpr:
			if newName, ok := funcRenameMap[v.Sel.Name]; ok {
				fmt.Printf("  ✏️  Updating method call: %s → %s in %s\n", v.Sel.Name, newName, path)
				v.Sel.Name = newName
				updated = true
			}
		}
		return true
	})

	if updated {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		err = printer.Fprint(file, fs, node)
		if err == nil {
			fmt.Printf("✅ Updated file: %s\n", path)
		}
		return err
	}

	return nil
}

func toCamelCasePreserveFirst(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) < 2 {
		return s
	}

	isExported := isUpper(parts[0][0])
	for i := range parts {
		if i == 0 {
			if isExported {
				parts[i] = strings.Title(parts[i])
			} else {
				parts[i] = strings.ToLower(parts[i])
			}
		} else {
			parts[i] = strings.Title(parts[i])
		}
	}
	return strings.Join(parts, "")
}

func isUpper(ch byte) bool {
	return 'A' <= ch && ch <= 'Z'
}

