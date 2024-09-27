package main

import "math/big"

func isAlpha(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'z') || (c == '_')
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func stringToBigFloat(str string) big.Float {
	var power, digit, ten, tmp, float_literal big.Float
	power.SetFloat64(1.0)
	ten.SetFloat64(10.0)
	float_literal.SetFloat64(0)
	found_period := false
	for _, char := range str {
		if char == '.' {
			found_period = true
			continue
		}
		if found_period {
			digit.SetFloat64(float64(char - '0'))
			power.Quo(&power, &ten)
			tmp.Mul(&power, &digit)
			float_literal.Add(&float_literal, &tmp)
		} else {
			digit.SetFloat64(float64(char - '0'))
			tmp.Mul(&float_literal, &ten)
			float_literal.Add(&tmp, &digit)
		}
	}
	return float_literal
}
