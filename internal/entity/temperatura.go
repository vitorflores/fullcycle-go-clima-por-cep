package entity

import "math"

type Temperatura struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NovaTemperatura(celsius float64) Temperatura {
	return Temperatura{
		Celsius:    arredondar(celsius),
		Fahrenheit: CelsiusParaFahrenheit(celsius),
		Kelvin:     CelsiusParaKelvin(celsius),
	}
}

func CelsiusParaFahrenheit(celsius float64) float64 {
	return arredondar(celsius*1.8 + 32)
}

func CelsiusParaKelvin(celsius float64) float64 {
	return arredondar(celsius + 273)
}

func arredondar(valor float64) float64 {
	return math.Round(valor*100) / 100
}
