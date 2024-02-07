package romBuilder

type bankFormat struct {
	chrSize         int
	mapper          int
	mirroring       int
	prgSize         int
	hasSetChr       bool
	hasSetMapper    bool
	hasSetMirroring bool
	hasSetPrg       bool
}

var BankInfo = bankFormat{}
