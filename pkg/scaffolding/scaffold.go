package scaffolding

import (
	"fmt"
)

func Process(filePath string) {
	model := ParseModel(filePath)

	fmt.Printf("%+v", model)
}
