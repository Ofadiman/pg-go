package utils

import (
	"encoding/json"
	"fmt"
)

func Print(value any) {
	b, _ := json.MarshalIndent(value, "", "  ")
	fmt.Println(string(b))
}
