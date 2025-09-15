package htpasswd

import (
	"crypto/sha512"
)

// sha512Crypt implements the SHA-512 crypt algorithm as specified in
// http://www.akkadia.org/drepper/SHA-crypt.txt
func sha512Crypt(password, salt string) string {
	const rounds = 5000
	const prefix = "$6$"

	// Custom base64 alphabet used by crypt
	const alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Ensure salt is maximum 16 characters
	if len(salt) > 16 {
		salt = salt[:16]
	}

	// Step 1: Compute alternate sum
	h := sha512.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))
	h.Write([]byte(password))
	altResult := h.Sum(nil)

	// Step 2: Compute main sum
	h.Reset()
	h.Write([]byte(password))
	h.Write([]byte(salt))

	// Add altResult for each character in password
	for i := len(password); i > 0; i -= 64 {
		if i > 64 {
			h.Write(altResult)
		} else {
			h.Write(altResult[:i])
		}
	}

	// Add password or altResult based on password length bits
	for i := len(password); i > 0; i >>= 1 {
		if (i & 1) != 0 {
			h.Write(altResult)
		} else {
			h.Write([]byte(password))
		}
	}

	result := h.Sum(nil)

	// Step 3: Compute P sequence
	h.Reset()
	for i := 0; i < len(password); i++ {
		h.Write([]byte(password))
	}
	pBytes := h.Sum(nil)

	// Create P sequence
	p := make([]byte, 0, len(password))
	for i := len(password); i > 0; i -= 64 {
		if i > 64 {
			p = append(p, pBytes...)
		} else {
			p = append(p, pBytes[:i]...)
		}
	}

	// Step 4: Compute S sequence
	h.Reset()
	for i := 0; i < 16+int(result[0]); i++ {
		h.Write([]byte(salt))
	}
	sBytes := h.Sum(nil)

	// Create S sequence
	s := make([]byte, 0, len(salt))
	for i := len(salt); i > 0; i -= 64 {
		if i > 64 {
			s = append(s, sBytes...)
		} else {
			s = append(s, sBytes[:i]...)
		}
	}

	// Step 5: Perform rounds iterations
	for round := 0; round < rounds; round++ {
		h.Reset()

		if (round & 1) != 0 {
			h.Write(p)
		} else {
			h.Write(result)
		}

		if round%3 != 0 {
			h.Write(s)
		}

		if round%7 != 0 {
			h.Write(p)
		}

		if (round & 1) != 0 {
			h.Write(result)
		} else {
			h.Write(p)
		}

		result = h.Sum(nil)
	}

	// Step 6: Encode result using custom base64
	encoded := make([]byte, 86)

	// Rearrange bytes according to the algorithm
	indices := [][]int{
		{0, 21, 42}, {22, 43, 1}, {44, 2, 23}, {3, 24, 45},
		{25, 46, 4}, {47, 5, 26}, {6, 27, 48}, {28, 49, 7},
		{50, 8, 29}, {9, 30, 51}, {31, 52, 10}, {53, 11, 32},
		{12, 33, 54}, {34, 55, 13}, {56, 14, 35}, {15, 36, 57},
		{37, 58, 16}, {59, 17, 38}, {18, 39, 60}, {40, 61, 19},
		{62, 20, 41}, {63},
	}

	pos := 0
	for _, group := range indices {
		val := 0
		for i, idx := range group {
			val |= int(result[idx]) << (i * 8)
		}

		for i := 0; i < 4 && pos < len(encoded); i++ {
			encoded[pos] = alphabet[val&0x3f]
			val >>= 6
			pos++
		}
	}

	// Trim trailing characters and return formatted result
	return prefix + salt + "$" + string(encoded[:86])
}
