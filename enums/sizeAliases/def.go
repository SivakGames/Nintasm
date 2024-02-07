package enumSizeAliases

type Def int

const (
	None Def = iota + 0
	Size1kb
	Size2kb
	Size4kb
	Size8kb
	Size16kb
	Size32kb
	Size64kb
	Size128kb
	Size256kb
	Size512kb
	Size1mb
	Size2mb
)
