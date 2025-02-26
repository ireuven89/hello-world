package utils

type Secrets struct {
	Name    string
	Secrete []Secrets
}

func ListSecrets(root Secrets, prefix string) []string {
	var paths []string
	currentPath := prefix + root.Name

	// If it's a leaf node, add the path
	if len(root.Secrete) == 0 {
		return []string{currentPath}
	}

	// Recursively get paths for child secrets
	for _, s := range root.Secrete {
		paths = append(paths, ListSecrets(s, currentPath+"/")...)
	}
	return paths
}
