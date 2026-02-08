package lobster

import (
	"errors"
	"strings"
)

const (
	Carrier = "🦞"
	ZWSP    = '\u200B'
	ZWNJ    = '\u200C'
	ZWJ     = '\u200D'
)

func EncodeTernary(data []byte) string {
	var sb strings.Builder
	sb.WriteString(Carrier)

	for _, b := range data {
		val := int(b)
		trits := make([]rune, 6)
		for i := 0; i < 6; i++ {
			rem := val % 3
			val /= 3
			switch rem {
			case 0:
				trits[i] = ZWSP
			case 1:
				trits[i] = ZWNJ
			case 2:
				trits[i] = ZWJ
			}
		}
		for _, t := range trits {
			sb.WriteRune(t)
		}
	}
	sb.WriteString(Carrier)
	return sb.String()
}

func EncodeBinary(data []byte, key string) string {
	var sb strings.Builder
	sb.WriteString(Carrier)

	keyBytes := []byte(key)
	keyLen := len(keyBytes)

	for i, b := range data {
		currentByte := b
		if keyLen > 0 {
			currentByte = b ^ keyBytes[i%keyLen]
		}

		for bit := 0; bit < 8; bit++ {
			if (currentByte>>bit)&1 == 1 {
				sb.WriteRune(ZWNJ)
			} else {
				sb.WriteRune(ZWSP)
			}
		}
	}
	sb.WriteString(Carrier)
	return sb.String()
}

// Decode
func Decode(encoded string, version int, key string) (string, error) {
	clean := strings.Trim(encoded, Carrier)

	runes := []rune(clean)
	var result []byte

	if version <= 3 {
		if len(runes)%6 != 0 {
			return "", errors.New("invalid ternary length")
		}
		for i := 0; i < len(runes); i += 6 {
			var val int
			multiplier := 1
			for j := 0; j < 6; j++ {
				r := runes[i+j]
				digit := 0
				if r == ZWNJ {
					digit = 1
				}
				if r == ZWJ {
					digit = 2
				}
				val += digit * multiplier
				multiplier *= 3
			}
			result = append(result, byte(val))
		}
	} else {
		if len(runes)%8 != 0 {
			return "", errors.New("invalid binary length")
		}
		for i := 0; i < len(runes); i += 8 {
			var val byte
			for j := 0; j < 8; j++ {
				if runes[i+j] == ZWNJ {
					val |= (1 << j)
				}
			}
			result = append(result, val)
		}

		if len(key) > 0 {
			keyBytes := []byte(key)
			for i := 0; i < len(result); i++ {
				result[i] ^= keyBytes[i%len(keyBytes)]
			}
		}
	}

	return string(result), nil
}
