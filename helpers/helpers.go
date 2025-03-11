package helpers

import (
	// Go Internal Packages
	"encoding/json"
	"fmt"
)

// PrintStruct prints a givens struct in pretty format with indent
func PrintStruct(v any) {
	res, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(res))
}
