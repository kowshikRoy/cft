package model

type Language struct {
	// Name of the Language
	Name string

	// Compiler Used in
	Compiler map[string]string

	// Extension of file for the language
	Extension string

	// Build flags to compile and build the source code
	BuildFlags string

	// Output file extension
	OutputFileExtension map[string]string
}
