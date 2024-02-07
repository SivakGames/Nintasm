package util

import (
	enumSizeAliases "misc/nintasm/enums/sizeAliases"
	"strings"
)

var sizeStringAliases = map[string]enumSizeAliases.Def{
	"1kb":   enumSizeAliases.Size1kb,
	"2kb":   enumSizeAliases.Size2kb,
	"4kb":   enumSizeAliases.Size4kb,
	"8kb":   enumSizeAliases.Size8kb,
	"16kb":  enumSizeAliases.Size16kb,
	"32kb":  enumSizeAliases.Size32kb,
	"64kb":  enumSizeAliases.Size64kb,
	"128kb": enumSizeAliases.Size128kb,
	"256kb": enumSizeAliases.Size256kb,
	"512kb": enumSizeAliases.Size512kb,
	"1mb":   enumSizeAliases.Size1mb,
	"2mb":   enumSizeAliases.Size2mb,
}

var sizeNumericAliases = map[int]enumSizeAliases.Def{
	0x0000400: enumSizeAliases.Size1kb,
	0x0000800: enumSizeAliases.Size2kb,
	0x0001000: enumSizeAliases.Size4kb,
	0x0002000: enumSizeAliases.Size8kb,
	0x0004000: enumSizeAliases.Size16kb,
	0x0008000: enumSizeAliases.Size32kb,
	0x0010000: enumSizeAliases.Size64kb,
	0x0020000: enumSizeAliases.Size128kb,
	0x0040000: enumSizeAliases.Size256kb,
	0x0080000: enumSizeAliases.Size512kb,
	0x0100000: enumSizeAliases.Size1mb,
	0x0200000: enumSizeAliases.Size2mb,
}

func ValidateSizeStringAlias(sizeAlias string) (enumSizeAliases.Def, bool) {
	adjustedAlias := strings.ToLower(sizeAlias)
	enumValue, enumOk := sizeStringAliases[adjustedAlias]
	if !enumOk {
		return enumSizeAliases.None, enumOk
	}
	return enumValue, enumOk
}

func ValidateSizeNumericAlias(sizeAlias int) (enumSizeAliases.Def, bool) {
	enumValue, enumOk := sizeNumericAliases[sizeAlias]
	if !enumOk {
		return enumSizeAliases.None, enumOk
	}
	return enumValue, enumOk
}
