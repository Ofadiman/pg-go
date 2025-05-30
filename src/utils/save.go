package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPackageName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("failed to get caller information"))
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, "/")
	if len(parts) == 0 {
		panic("parts are empty")
	}

	lastDot := strings.LastIndexByte(parts[len(parts)-1], '.')
	if lastDot < 0 {
		panic("dot not found")
	}

	return parts[len(parts)-1][:lastDot]
}

func Save(packageName string, data any) {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(fmt.Errorf("failed to get caller information"))
	}

	funcName := runtime.FuncForPC(pc).Name()
	if lastDot := strings.LastIndex(funcName, "."); lastDot != -1 {
		funcName = funcName[lastDot+1:]
	}

	sourceFile := filepath.Base(file)
	sourceFile = strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile))

	filename := fmt.Sprintf("tmp/%s_%s.json", packageName, funcName)

	if err := os.MkdirAll("tmp", 0755); err != nil {
		panic(fmt.Errorf("failed to create tmp directory: %w", err))
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal JSON: %w", err))
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		panic(fmt.Errorf("failed to write file %s: %w", filename, err))
	}

	fmt.Printf("successfully saved json to: %s\n", filename)
}
