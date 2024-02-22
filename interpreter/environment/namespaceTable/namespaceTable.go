package namespaceTable

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

var namespaceSymbolTable = map[string]bool{}

func AddNamespaceToEnvironment(namespaceName string) {
	namespaceSymbolTable[namespaceName] = true
}
