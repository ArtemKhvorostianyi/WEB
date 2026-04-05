package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Result struct {
	Pc, Price, Sigma1, Sigma2 float64

	A, B float64

	W1, W2 float64
	W3, W4 float64

	EnergyTotal float64

	P1, S1 float64
	P2, S2 float64

	Profit float64

	Show bool
}

func parseFloat(r *http.Request, key string) float64 {
	val, _ := strconv.ParseFloat(r.FormValue(key), 64)
	return val
}

// нормальний розподіл
func normalPDF(x, mean, sigma float64) float64 {
	return (1 / (sigma * math.Sqrt(2*math.Pi))) *
		math.Exp(-math.Pow(x-mean, 2)/(2*sigma*sigma))
}

// чисельний інтеграл
func integrate(mean, sigma, a, b float64) float64 {
	steps := 10000.0
	h := (b - a) / steps
	sum := 0.0

	for i := 0.0; i < steps; i++ {
		x := a + i*h
		sum += normalPDF(x, mean, sigma) * h
	}

	return sum
}

func handler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("Практична 3.html"))
	data := Result{
		Pc: 5, Price: 7, Sigma1: 1, Sigma2: 0.25,
	}

	if r.Method == http.MethodPost {

		data.Pc = parseFloat(r, "pc")
		data.Price = parseFloat(r, "price")
		data.Sigma1 = parseFloat(r, "sigma1")
		data.Sigma2 = parseFloat(r, "sigma2")

		data.A = data.Pc * 0.95
		data.B = data.Pc * 1.05

		w1 := integrate(data.Pc, data.Sigma1, data.A, data.B)
		w2 := integrate(data.Pc, data.Sigma2, data.A, data.B)

		data.EnergyTotal = data.Pc * 24

		data.W1 = data.EnergyTotal * w1
		data.W2 = data.EnergyTotal * (1 - w1)

		data.P1 = data.W1 * data.Price
		data.S1 = data.W2 * data.Price

		data.W3 = data.EnergyTotal * w2
		data.W4 = data.EnergyTotal * (1 - w2)

		data.P2 = data.W3 * data.Price
		data.S2 = data.W4 * data.Price

		data.Profit = data.P2 - data.S2

		data.Show = true
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ПЗ3: http://localhost:8082")
	http.ListenAndServe(":8082", nil)
}