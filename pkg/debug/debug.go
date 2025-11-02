package debug

import (
	"encoding/json"
	"fmt"
)

func PrintStructJson(strt interface{}) {
	jsonData, err := json.MarshalIndent(strt, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}
