package namespaceTable

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

var IsDefiningNamespace bool

type namespaceEntry struct {
	Key      string
	Resolved bool
}

func newNamespaceEntry(key string, resolved bool) namespaceEntry {
	return namespaceEntry{
		Key:      key,
		Resolved: resolved,
	}
}

// ============================================================

// Mainly just acts a true/false entry as to whether the NS exists.
// Actual values are stored in the regular symbol table.
var namespaceSymbolTable = map[string][]namespaceEntry{}

func AddIdentifierKeyToNamespaceTable(namespaceName string) {
	namespaceSymbolTable[namespaceName] = []namespaceEntry{}
}

func AddKeyToCurrentNamespace(namespaceName string, key string, resolved bool) {
	namespaceSymbolTable[namespaceName] = append(
		namespaceSymbolTable[namespaceName],
		newNamespaceEntry(key, resolved))
}

func GetNamespaceValues(namespaceName string) (*[]namespaceEntry, error) {
	entries, ok := namespaceSymbolTable[namespaceName]
	if !ok {
		return nil, errorHandler.AddNew(enumErrorCodes.NamespaceNotExist, namespaceName)
	}
	return &entries, nil
}
