package htpasswd

import (
	"strings"
	"testing"
)

func TestSha512Crypt(t *testing.T) {
	password := "secret123"
	salt := "saltySal"

	result := sha512Crypt(password, salt)

	// Check that result has the correct prefix
	expectedPrefix := "$6$" + salt + "$"
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected prefix %s, but got %s", expectedPrefix, result[:len(expectedPrefix)])
	}

	// Check that result has reasonable length (should be around 98-106 characters)
	if len(result) < 90 || len(result) > 110 {
		t.Errorf("Expected result length between 90-110, got %d", len(result))
	}

	// Test with empty salt
	result2 := sha512Crypt(password, "")
	if !strings.HasPrefix(result2, "$6$$") {
		t.Errorf("Expected empty salt to work, got %s", result2)
	}

	// Test that same input produces same output
	result3 := sha512Crypt(password, salt)
	if result != result3 {
		t.Errorf("Expected deterministic output, got different results")
	}
}
