package options

// Count returns the number of options in the catalog.
func Count() int {
	n := 0
	for _, sec := range Catalog {
		n += len(sec.Items)
	}
	return n
}
