package utils

import (
	"encoding/json"
	"fmt"
)

func Print(value any) {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
