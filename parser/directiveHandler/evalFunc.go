package directiveHandler

import "fmt"

func evalFunc(operandList *[]Node) error {
	functionNode := (*operandList)[0]
	fmt.Println(functionNode)
	return nil
}
