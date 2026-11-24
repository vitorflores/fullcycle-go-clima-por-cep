package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCEPValido(t *testing.T) {
	tabela := []struct {
		cep      string
		esperado bool
	}{
		{"01001000", true},
		{"29902555", true},
		{"0100100", false},
		{"010010000", false},
		{"01001-000", false},
		{"0100100a", false},
		{"", false},
		{"        ", false},
	}
	for _, item := range tabela {
		assert.Equal(t, item.esperado, CEPValido(item.cep), "cep %q", item.cep)
	}
}
