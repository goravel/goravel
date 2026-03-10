package maputil

import "fmt"

// Aliases implemented a simple string alias map.
type Aliases map[string]string

// AddAlias to the Aliases
func (as Aliases) AddAlias(alias, real string) {
	if rn, ok := as[alias]; ok {
		panic(fmt.Sprintf("The alias '%s' is already used by '%s'", alias, rn))
	}
	as[alias] = real
}

// AddAliases to the Aliases
func (as Aliases) AddAliases(real string, aliases []string) {
	for _, a := range aliases {
		as.AddAlias(a, real)
	}
}

// AddAliasMap to the Aliases
func (as Aliases) AddAliasMap(alias2real map[string]string) {
	for a, r := range alias2real {
		as.AddAlias(a, r)
	}
}

// HasAlias in the Aliases
func (as Aliases) HasAlias(alias string) bool {
	if _, ok := as[alias]; ok {
		return true
	}
	return false
}

// ResolveAlias by given name.
func (as Aliases) ResolveAlias(alias string) string {
	if name, ok := as[alias]; ok {
		return name
	}
	return alias
}
