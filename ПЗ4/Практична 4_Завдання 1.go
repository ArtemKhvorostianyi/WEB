package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type Result struct {
	U, Sm, Ik, T, Jek, C float64

	Im  float64
	Ipa float64

	Sek  float64
	Smin float64

	Cable float64

	Show bool
}

func parseFloat(r *http.Request, key string) float64 {
	val, _ := strconv.ParseFloat(r.FormValue(key), 64)
	return val
}

func handler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("Практична 4_Завдання 1.html"))

	data := Result{
		U: 10, Sm: 1300, Ik: 2.5, T: 2.5, Jek: 1.4, C: 92,
	}

	if r.Method == http.MethodPost {

		data.U = parseFloat(r, "voltage")
		data.Sm = parseFloat(r, "sm")
		data.Ik = parseFloat(r, "ik")
		data.T = parseFloat(r, "time")
		data.Jek = parseFloat(r, "jek")
		data.C = parseFloat(r, "c")

		data.Im = (data.Sm / 2) / (math.Sqrt(3) * data.U)
		data.Ipa = 2 * data.Im

		data.Sek = data.Im / data.Jek
		data.Smin = (data.Ik * 1000 * math.Sqrt(data.T)) / data.C

		standard := []float64{16, 25, 35, 50, 70, 95, 120}

		required := math.Max(data.Sek, data.Smin)

		for _, s := range standard {
			if s >= required {
				data.Cable = s
				break
			}
		}

		data.Show = true
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ПЗ4: http://localhost:8083")
	http.ListenAndServe(":8083", nil)
}