package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Result struct {
	Coal, Mazut, Gas float64

	KCoal, KMazut float64
	ECoal, EMazut float64

	Show bool
}

func parseFloat(r *http.Request, key string) float64 {
	val, _ := strconv.ParseFloat(r.FormValue(key), 64)
	return val
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("Практична 2.html"))

	data := Result{}

	if r.Method == http.MethodPost {

		data.Coal = parseFloat(r, "coal")
		data.Mazut = parseFloat(r, "mazut")
		data.Gas = parseFloat(r, "gas")

		eta := 0.985

		Qcoal := 20.47
		Acoal := 25.20
		avin := 0.8
		G := 1.5

		Qmazut := 39.48
		Amazut := 0.15

		kCoal := (1000000 / Qcoal) * avin * (Acoal / (100 - G)) * (1 - eta)
		Ecoal := 1e-6 * kCoal * Qcoal * data.Coal * 1000

		kMazut := (1000000 / Qmazut) * (Amazut / 100) * (1 - eta)
		Emazut := 1e-6 * kMazut * Qmazut * data.Mazut * 1000

		data.KCoal = kCoal
		data.KMazut = kMazut
		data.ECoal = Ecoal
		data.EMazut = Emazut

		data.Show = true
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ПЗ2: http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}