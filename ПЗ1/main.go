package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Result struct {
	Error string

	Krs float64
	Krg float64

	Hc, Cc, Sc, Nc, Oc, Ac float64
	Hg, Cg, Sg, Ng, Og     float64

	Qr, Qd, Qdaf float64
	Show         bool
}

func parseFloat(r *http.Request, key string) float64 {
	val, _ := strconv.ParseFloat(r.FormValue(key), 64)
	return val
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("_ПВЗ_Task_1.html"))

	data := Result{}

	if r.Method == http.MethodPost {
		Hp := parseFloat(r, "Hp")
		Cp := parseFloat(r, "Cp")
		Sp := parseFloat(r, "Sp")
		Np := parseFloat(r, "Np")
		Op := parseFloat(r, "Op")
		Wp := parseFloat(r, "Wp")
		Ap := parseFloat(r, "Ap")

		sum := Hp + Cp + Sp + Np + Op + Wp + Ap

		if abs(sum-100) > 0.01 {
			data.Error = "❌ Сума компонентів повинна дорівнювати 100%"
		} else {
			Krs := 100 / (100 - Wp)
			Krg := 100 / (100 - Wp - Ap)

			data.Krs = Krs
			data.Krg = Krg

			data.Hc = Hp * Krs
			data.Cc = Cp * Krs
			data.Sc = Sp * Krs
			data.Nc = Np * Krs
			data.Oc = Op * Krs
			data.Ac = Ap * Krs

			data.Hg = Hp * Krg
			data.Cg = Cp * Krg
			data.Sg = Sp * Krg
			data.Ng = Np * Krg
			data.Og = Op * Krg

			Qr_kJ := 339*Cp + 1030*Hp - 108.8*(Op-Sp) - 25*Wp
			Qr := Qr_kJ / 1000
			Qd := (Qr + 0.025*Wp) * 100 / (100 - Wp)
			Qdaf := (Qr + 0.025*Wp) * 100 / (100 - Wp - Ap)

			data.Qr = Qr
			data.Qd = Qd
			data.Qdaf = Qdaf

			data.Show = true
		}
	}

	tmpl.Execute(w, data)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("ПЗ1: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}