package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// дані одного електроприймача
type EP struct {
	Name string
	N    int
	Pn   float64
	Kv   float64
	Tgf  float64
	Eta  float64
	Cosf float64
	U    float64
}

func (e EP) NPn() float64     { return float64(e.N) * e.Pn }
func (e EP) NPnKv() float64   { return e.NPn() * e.Kv }
func (e EP) NPnKvTg() float64 { return e.NPnKv() * e.Tgf }
func (e EP) NPn2() float64    { return float64(e.N) * e.Pn * e.Pn }

// розрахунковий струм ЕП (п.3.2)
func (e EP) Ip() float64 {
	zn := math.Sqrt(3) * e.U * e.Cosf * e.Eta
	if zn == 0 {
		return 0
	}
	return e.NPn() / zn
}

// результат розрахунку групи
type GroupResult struct {
	SumNPn     float64
	SumNPnKv   float64
	SumNPnKvTg float64
	SumNPn2    float64
	TotalN     int
	KvGroup    float64
	Ne         int
	NeRaw      float64
	Kp         float64
	Pp         float64
	Qp         float64
	Sp         float64
	IpGroup    float64
	U          float64
	TableName  string
}

// дані для шаблону
type PageData struct {
	ShrList     []EP
	BigList     []EP
	ShrCount    int
	ResSHR      *GroupResult
	ResWorkshop *GroupResult
	Calculated  bool
}

// таблиця 6.3 (Кр, T0=10хв) 

var neRows63 = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	12, 14, 16, 18, 20, 25, 30, 35, 40, 50, 60, 80, 100,
}

var kvCols63 = []float64{0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

var table63 = [][]float64{
	{8.00, 5.33, 4.00, 2.67, 2.00, 1.60, 1.33, 1.14, 1.00},
	{6.22, 4.33, 3.39, 2.45, 1.98, 1.60, 1.33, 1.14, 1.00},
	{4.06, 2.89, 2.31, 1.74, 1.45, 1.34, 1.22, 1.14, 1.00},
	{3.24, 2.35, 1.91, 1.47, 1.25, 1.21, 1.12, 1.06, 1.00},
	{2.84, 2.09, 1.72, 1.35, 1.16, 1.16, 1.08, 1.03, 1.00},
	{2.64, 1.96, 1.62, 1.28, 1.14, 1.13, 1.06, 1.01, 1.00},
	{2.49, 1.86, 1.54, 1.23, 1.12, 1.10, 1.04, 1.00, 1.00},
	{2.37, 1.78, 1.48, 1.19, 1.10, 1.08, 1.02, 1.00, 1.00},
	{2.27, 1.71, 1.43, 1.16, 1.09, 1.07, 1.01, 1.00, 1.00},
	{2.18, 1.65, 1.39, 1.13, 1.07, 1.05, 1.00, 1.00, 1.00},
	{2.04, 1.56, 1.32, 1.08, 1.05, 1.03, 1.00, 1.00, 1.00},
	{1.94, 1.49, 1.27, 1.05, 1.02, 1.00, 1.00, 1.00, 1.00},
	{1.85, 1.43, 1.23, 1.02, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.78, 1.39, 1.19, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.72, 1.35, 1.16, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.60, 1.27, 1.10, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.51, 1.21, 1.05, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.44, 1.16, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.40, 1.13, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.30, 1.07, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.25, 1.03, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.16, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
	{1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00, 1.00},
}

// таблиця 6.4 (Кр, T0=2.5год)

var kvCols64 = []float64{0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7}

var table64 = [][]float64{
	{8.00, 5.33, 4.00, 2.67, 2.00, 1.60, 1.33, 1.14},
	{5.01, 3.44, 2.69, 1.90, 1.52, 1.24, 1.11, 1.00},
	{2.40, 2.17, 1.80, 1.42, 1.23, 1.14, 1.08, 1.00},
	{2.28, 1.73, 1.46, 1.19, 1.06, 1.04, 1.00, 0.97},
	{1.31, 1.12, 1.02, 1.00, 0.98, 0.96, 0.94, 0.93},
	{1.20, 1.00, 0.96, 0.95, 0.94, 0.93, 0.92, 0.91},
	{1.10, 0.97, 0.91, 0.90, 0.90, 0.90, 0.90, 0.90},
	{0.80, 0.80, 0.80, 0.85, 0.85, 0.85, 0.90, 0.90},
	{0.75, 0.75, 0.75, 0.75, 0.75, 0.80, 0.85, 0.85},
	{0.65, 0.65, 0.65, 0.70, 0.70, 0.75, 0.80, 0.80},
}

// пошук найближчого стовпця Кв
func findKvCol(cols []float64, kv float64) int {
	idx := 0
	minD := math.Abs(cols[0] - kv)
	for i := 1; i < len(cols); i++ {
		d := math.Abs(cols[i] - kv)
		if d < minD {
			minD = d
			idx = i
		}
	}
	return idx
}

// знайти Кр по табл 6.3
func getKp63(kv float64, ne int) float64 {
	rowIdx := 0
	for i := len(neRows63) - 1; i >= 0; i-- {
		if neRows63[i] <= ne {
			rowIdx = i
			break
		}
	}
	colIdx := findKvCol(kvCols63, kv)
	return table63[rowIdx][colIdx]
}

// знайти Кр по табл 6.4
func getKp64(kv float64, ne int) float64 {
	var rowIdx int
	switch {
	case ne <= 1:
		rowIdx = 0
	case ne <= 2:
		rowIdx = 1
	case ne <= 3:
		rowIdx = 2
	case ne <= 4:
		rowIdx = 3
	case ne <= 5:
		rowIdx = 4
	case ne <= 8:
		rowIdx = 5
	case ne <= 10:
		rowIdx = 6
	case ne <= 25:
		rowIdx = 7
	case ne <= 50:
		rowIdx = 8
	default:
		rowIdx = 9
	}

	colIdx := findKvCol(kvCols64, kv)
	if kv >= 0.7 {
		colIdx = len(kvCols64) - 1
	}

	return table64[rowIdx][colIdx]
}

// розрахунок для групи (ШР) - рівень II
func calcSHR(list []EP) GroupResult {
	r := GroupResult{TableName: "6.3"}
	if len(list) == 0 {
		return r
	}
	r.U = list[0].U

	for _, ep := range list {
		r.SumNPn += ep.NPn()
		r.SumNPnKv += ep.NPnKv()
		r.SumNPnKvTg += ep.NPnKvTg()
		r.SumNPn2 += ep.NPn2()
		r.TotalN += ep.N
	}

	// п.4.1 груповий Кв
	if r.SumNPn > 0 {
		r.KvGroup = r.SumNPnKv / r.SumNPn
	}

	// п.4.2 ne (округлення до меншого цілого)
	if r.SumNPn2 > 0 {
		r.NeRaw = (r.SumNPn * r.SumNPn) / r.SumNPn2
		r.Ne = int(r.NeRaw)
		if r.Ne < 1 {
			r.Ne = 1
		}
	}

	// п.4.3 Кр з табл 6.3
	r.Kp = getKp63(r.KvGroup, r.Ne)

	// п.4.4 Рр
	r.Pp = r.Kp * r.SumNPnKv

	// п.4.5 Qр (формула 6.8)
	if r.Ne <= 10 {
		r.Qp = 1.1 * r.SumNPnKvTg
	} else {
		r.Qp = r.SumNPnKvTg
	}

	// п.4.6 Sр
	r.Sp = math.Sqrt(r.Pp*r.Pp + r.Qp*r.Qp)

	// п.4.7 Ір = Sр / (√3 * Uн)
	if r.U > 0 {
		r.IpGroup = r.Sp / (math.Sqrt(3) * r.U)
	}

	return r
}

// розрахунок цеху (рівень III)
func calcWorkshop(shrList []EP, shrCount int, bigList []EP) GroupResult {
	r := GroupResult{TableName: "6.4"}

	if len(shrList) > 0 {
		r.U = shrList[0].U
	} else if len(bigList) > 0 {
		r.U = bigList[0].U
	}

	// сумуємо по ШР (з множенням на к-ть ідентичних)
	for _, ep := range shrList {
		r.SumNPn += ep.NPn() * float64(shrCount)
		r.SumNPnKv += ep.NPnKv() * float64(shrCount)
		r.SumNPnKvTg += ep.NPnKvTg() * float64(shrCount)
		r.SumNPn2 += ep.NPn2() * float64(shrCount)
		r.TotalN += ep.N * shrCount
	}

	// додаємо крупні ЕП
	for _, ep := range bigList {
		r.SumNPn += ep.NPn()
		r.SumNPnKv += ep.NPnKv()
		r.SumNPnKvTg += ep.NPnKvTg()
		r.SumNPn2 += ep.NPn2()
		r.TotalN += ep.N
	}

	// п.6.1 Кв цеху
	if r.SumNPn > 0 {
		r.KvGroup = r.SumNPnKv / r.SumNPn
	}

	// п.6.2 ne цеху
	if r.SumNPn2 > 0 {
		r.NeRaw = (r.SumNPn * r.SumNPn) / r.SumNPn2
		r.Ne = int(r.NeRaw)
		if r.Ne < 1 {
			r.Ne = 1
		}
	}

	// п.6.3 Кр з табл 6.4
	r.Kp = getKp64(r.KvGroup, r.Ne)

	// п.6.4 Рр
	r.Pp = r.Kp * r.SumNPnKv

	// п.6.5 Qр (формула 6.8 - без Кр!)
	if r.Ne <= 10 {
		r.Qp = 1.1 * r.SumNPnKvTg
	} else {
		r.Qp = r.SumNPnKvTg
	}

	// п.6.6 Sр
	r.Sp = math.Sqrt(r.Pp*r.Pp + r.Qp*r.Qp)

	// п.6.7 Ір
	if r.U > 0 {
		r.IpGroup = r.Sp / (math.Sqrt(3) * r.U)
	}

	return r
}

// контрольний приклад
func defaultSHR() []EP {
	return []EP{
		{Name: "Шліфувальний верстат (1-4)", N: 4, Pn: 20, Kv: 0.15, Tgf: 1.33, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Свердлильний верстат (5-6)", N: 2, Pn: 14, Kv: 0.12, Tgf: 1.0, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Фугувальний верстат (9-12)", N: 4, Pn: 42, Kv: 0.15, Tgf: 1.33, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Циркулярна пила (13)", N: 1, Pn: 36, Kv: 0.3, Tgf: 1.52, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Прес (16)", N: 1, Pn: 20, Kv: 0.5, Tgf: 0.75, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Полірувальний верстат (24)", N: 1, Pn: 40, Kv: 0.2, Tgf: 1.0, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Фрезерний верстат (26-27)", N: 2, Pn: 32, Kv: 0.2, Tgf: 1.0, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Вентилятор (36)", N: 1, Pn: 20, Kv: 0.65, Tgf: 0.75, Eta: 0.92, Cosf: 0.9, U: 0.38},
	}
}

func defaultBigEP() []EP {
	return []EP{
		{Name: "Зварювальний трансформатор", N: 2, Pn: 100, Kv: 0.2, Tgf: 3.0, Eta: 0.92, Cosf: 0.9, U: 0.38},
		{Name: "Сушильна шафа", N: 2, Pn: 120, Kv: 0.8, Tgf: 0.0, Eta: 0.92, Cosf: 1.0, U: 0.38},
	}
}

// допоміжні функції для шаблону
func f2(v float64) string { return fmt.Sprintf("%.2f", v) }
func f1(v float64) string { return fmt.Sprintf("%.1f", v) }
func f0(v float64) string { return fmt.Sprintf("%.0f", v) }
func f4(v float64) string { return fmt.Sprintf("%.4f", v) }

// парсинг float з форми
func parseF(s string) float64 {
	s = strings.TrimSpace(s)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseI(s string) int {
	s = strings.TrimSpace(s)
	v, _ := strconv.Atoi(s)
	return v
}

// шаблон
var tpl *template.Template

func main() {
	funcMap := template.FuncMap{
		"f0": f0, "f1": f1, "f2": f2, "f4": f4,
	}
	tpl = template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html"))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/calc", handleCalc)

	log.Println("ПЗ6: http://localhost:8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ShrList:  defaultSHR(),
		BigList:  defaultBigEP(),
		ShrCount: 3,
	}
	tpl.Execute(w, data)
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	shrCount := parseI(r.FormValue("shrCount"))
	if shrCount < 1 {
		shrCount = 3
	}

	// парсимо ЕП з форми
	shrList := parseEPList(r, "shr")
	bigList := parseEPList(r, "big")

	resSHR := calcSHR(shrList)
	resWorkshop := calcWorkshop(shrList, shrCount, bigList)

	data := PageData{
		ShrList:     shrList,
		BigList:     bigList,
		ShrCount:    shrCount,
		ResSHR:      &resSHR,
		ResWorkshop: &resWorkshop,
		Calculated:  true,
	}
	tpl.Execute(w, data)
}

func parseEPList(r *http.Request, prefix string) []EP {
	var list []EP
	// шукаємо поки є name_0, name_1, ...
	for i := 0; i < 50; i++ {
		name := r.FormValue(fmt.Sprintf("%s_name_%d", prefix, i))
		if name == "" {
			continue
		}
		ep := EP{
			Name: name,
			N:    parseI(r.FormValue(fmt.Sprintf("%s_n_%d", prefix, i))),
			Pn:   parseF(r.FormValue(fmt.Sprintf("%s_pn_%d", prefix, i))),
			Kv:   parseF(r.FormValue(fmt.Sprintf("%s_kv_%d", prefix, i))),
			Tgf:  parseF(r.FormValue(fmt.Sprintf("%s_tgf_%d", prefix, i))),
			Eta:  parseF(r.FormValue(fmt.Sprintf("%s_eta_%d", prefix, i))),
			Cosf: parseF(r.FormValue(fmt.Sprintf("%s_cosf_%d", prefix, i))),
			U:    parseF(r.FormValue(fmt.Sprintf("%s_u_%d", prefix, i))),
		}
		if ep.Eta == 0 {
			ep.Eta = 0.92
		}
		if ep.Cosf == 0 {
			ep.Cosf = 0.9
		}
		if ep.U == 0 {
			ep.U = 0.38
		}
		list = append(list, ep)
	}
	return list
}
