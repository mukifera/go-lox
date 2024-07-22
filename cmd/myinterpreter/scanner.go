package main

func tokenize(lox_file_contents string) []string {
	var tokenization []string
	for i := 0; i < len(lox_file_contents); i++ {
		char := lox_file_contents[i]
		if char == '('{
			tokenization = append(tokenization, "LEFT_PAREN ( null")
		}else if char == ')' {
			tokenization = append(tokenization, "RIGHT_PAREN ) null")
		}else if char == '{' {
			tokenization = append(tokenization, "LEFT_BRACE { null")
		}else if char == '}' {
			tokenization = append(tokenization, "RIGHT_BRACE } null")
		}
	}
	tokenization = append(tokenization, "EOF  null")
	return tokenization
}