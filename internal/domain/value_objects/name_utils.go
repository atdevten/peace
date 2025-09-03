package value_objects

// isValidNameCharacter checks if a character is valid for a name
func isValidNameCharacter(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == ' ' ||
		char == '-' ||
		char == '\'' ||
		char == '.'
}
