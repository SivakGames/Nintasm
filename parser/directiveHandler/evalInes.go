package directiveHandler

import "misc/nintasm/romBuilder"

// +++++++++++++++++++++++++

func evalInesOperands(directiveName string, operandList *[]Node) error {
	inesNode := &(*operandList)[0]

	switch directiveName {
	case "INESPRG":
		return romBuilder.ValidateInesPrg(inesNode)
	case "INESCHR":
		return romBuilder.ValidateInesChr(inesNode)
	case "INESMAP":
		return romBuilder.ValidateInesMap(inesNode)
	case "INESMIR":
		return romBuilder.ValidateInesMirroring(inesNode)
	case "INESBAT":
		return nil
	default:
		panic("🛑 Something is very wrong with ines directive")
	}

}
