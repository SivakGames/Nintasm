package namespaceTable

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

// Mainly just acts a true/false entry as to whether the NS exists.
// Actual values are stored in the regular symbol table.
var namespaceSymbolTable = map[string]bool{}

func AddIdentifierKeyToNamespaceTable(namespaceName string) {
	namespaceSymbolTable[namespaceName] = true
}

func LookupInNamespaceTable(namespaceName string) {

}
