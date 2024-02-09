package handlerDirective

import "misc/nintasm/romBuilder"

// +++++++++++++++++++++++++

func evalInesOperands(directiveName string, operandList *[]Node) error {
	var err error
	inesNode := &(*operandList)[0]

	switch directiveName {
	case "INESPRG":
		err = romBuilder.ValidateInesPrg(inesNode)
	case "INESCHR":
		err = romBuilder.ValidateInesChr(inesNode)
	case "INESMAP":
		err = romBuilder.ValidateInesMap(inesNode)
	case "INESMIR":
		err = romBuilder.ValidateInesMirroring(inesNode)
	case "INESBAT":
		err = nil
	default:
		panic("Something is very wrong with ines directive")
	}

	return err
}
