package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusParaFahrenheit(t *testing.T) {
	tabela := []struct {
		celsius, esperado float64
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{28.5, 83.3},
		{36.6, 97.88},
	}
	for _, item := range tabela {
		assert.Equal(t, item.esperado, CelsiusParaFahrenheit(item.celsius), "celsius %v", item.celsius)
	}
}

func TestCelsiusParaKelvin(t *testing.T) {
	tabela := []struct {
		celsius, esperado float64
	}{
		{0, 273},
		{100, 373},
		{-273, 0},
		{28.5, 301.5},
	}
	for _, item := range tabela {
		assert.Equal(t, item.esperado, CelsiusParaKelvin(item.celsius), "celsius %v", item.celsius)
	}
}

func TestNovaTemperatura(t *testing.T) {
	temperatura := NovaTemperatura(28.5)
	assert.Equal(t, 28.5, temperatura.Celsius)
	assert.Equal(t, 83.3, temperatura.Fahrenheit)
	assert.Equal(t, 301.5, temperatura.Kelvin)
}

func FuzzCelsiusParaKelvin(f *testing.F) {
	for _, celsius := range []float64{-273, -10.5, 0, 25, 1000} {
		f.Add(celsius)
	}
	f.Fuzz(func(t *testing.T, celsius float64) {
		if celsius < -1e6 || celsius > 1e6 {
			return
		}
		resultado := CelsiusParaKelvin(celsius)
		if diferenca := resultado - (celsius + 273); diferenca > 0.005 || diferenca < -0.005 {
			t.Errorf("recebi %f para %f celsius", resultado, celsius)
		}
	})
}
