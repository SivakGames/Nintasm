package directiveHandler

import (
	"misc/nintasm/assemble/logHandler"
	"misc/nintasm/interpreter"
)

func evalLog(operandList *[]Node) error {
	logNode, err := interpreter.EvaluateNode((*operandList)[0])
	if err != nil {
		logHandler.AddLog("undefined on pass 1")
		return err
	}

	logHandler.AddLog(logNode.NodeValue)
	return nil
}
