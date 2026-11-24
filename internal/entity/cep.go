package entity

import (
	"errors"
	"regexp"
)

var (
	ErrCEPInvalido      = errors.New("invalid zipcode")
	ErrCEPNaoEncontrado = errors.New("can not find zipcode")
)

var formatoCEP = regexp.MustCompile(`^[0-9]{8}$`)

func CEPValido(cep string) bool {
	return formatoCEP.MatchString(cep)
}
