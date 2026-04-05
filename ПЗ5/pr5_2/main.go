package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Result1 struct {
	WOc, TBoc, KaOc, KpOc, WDk, WDs float64
}

type Result2 struct {
	MwA, MwP, TotalZ float64
}

type PageData struct {
	Tab  string
	Res1 *Result1
	Res2 *Result2
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "index.html", PageData{Tab: "task1"})
	})

	http.HandleFunc("/task1", calcTask1)
	http.HandleFunc("/task2", calcTask2)

	log.Println("ПЗ5: http://localhost:8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}

func calcTask1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w1 := parseF(r.FormValue("w1"))
	t1 := parseF(r.FormValue("t1"))
	w2 := parseF(r.FormValue("w2"))
	t2 := parseF(r.FormValue("t2"))
	w3 := parseF(r.FormValue("w3"))
	t3 := parseF(r.FormValue("t3"))
	w4 := parseF(r.FormValue("w4"))
	t4 := parseF(r.FormValue("t4"))
	w5 := parseF(r.FormValue("w5"))
	t5 := parseF(r.FormValue("t5"))
	kpMax := parseF(r.FormValue("kpMax"))
	wSection := parseF(r.FormValue("wSection"))

	wOc := w1 + w2 + w3 + w4 + w5
	tBoc := (w1*t1 + w2*t2 + w3*t3 + w4*t4 + w5*t5) / wOc
	kaOc := (wOc * tBoc) / 8760
	kpOc := 1.2 * (kpMax / 8760)

	wDk := 2 * wOc * (kaOc + kpOc)
	wDs := wDk + wSection

	res := &Result1{WOc: wOc, TBoc: tBoc, KaOc: kaOc, KpOc: kpOc, WDk: wDk, WDs: wDs}
	tpl.ExecuteTemplate(w, "index.html", PageData{Tab: "task1", Res1: res})
}

func calcTask2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	omega := parseF(r.FormValue("omega"))
	tB := parseF(r.FormValue("tB"))
	pm := parseF(r.FormValue("pm"))
	tm := parseF(r.FormValue("tm"))
	kp := parseF(r.FormValue("kp"))
	zA := parseF(r.FormValue("zA"))
	zP := parseF(r.FormValue("zP"))

	mwA := omega * tB * pm * tm
	mwP := kp * pm * tm
	totalZ := (zA * mwA) + (zP * mwP)

	res := &Result2{MwA: mwA, MwP: mwP, TotalZ: totalZ}
	tpl.ExecuteTemplate(w, "index.html", PageData{Tab: "task2", Res2: res})
}

func parseF(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
