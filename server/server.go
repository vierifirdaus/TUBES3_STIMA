package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	PertanyaanReq struct {
		Pertanyaan string `json:"pertanyaan"`
	}
	Pertanyaan struct {
		Pertanyaan string `json:"pertanyaan"`
		Jawaban    string `json:"jawaban"`
	}
	HistoriReq struct {
		Nama string `json:"nama"`
	}

	Histori struct {
		Nama string   `json:"nama"`
		Isi  []Respon `json:"isi"`
	}
	Respon struct {
		ID_histori int    `json:"id_histori"`
		Jenis      string `json:"jenis"`
		Isi        string `json:"isi"`
	}

	HistoriReqId struct {
		ID_histori int `json:"id_histori"`
	}

	HistoriId struct {
		ID_histori int    `json:"id"`
		Nama       string `json:"nama"`
	}
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL")+"/tubes3")
	if err != nil {
		fmt.Println("error")
	}
	return db, err
}

func getAllHistori(c echo.Context) error {
	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}
	defer db.Close()
	rows, err := db.Query("select * from histori")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer rows.Close()
	var isi []HistoriId
	for rows.Next() {
		var respon HistoriId
		err := rows.Scan(&respon.ID_histori, &respon.Nama)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}
		isi = append(isi, respon)
	}
	return c.JSON(http.StatusOK, isi)
}

func findAnswer(c echo.Context) error {
	var quest PertanyaanReq
	err := c.Bind(&quest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	quest.Pertanyaan = strings.ToLower(quest.Pertanyaan)

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	row := db.QueryRow("SELECT jawaban FROM pertanyaan WHERE LOWER(pertanyaan) = ?", quest.Pertanyaan)
	var answer string
	err = row.Scan(&answer)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusOK, answer)
}

func addQuestion(c echo.Context) error {
	var quest Pertanyaan
	err := c.Bind(&quest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	rows, err := db.Query("select Pertanyaan,Jawaban from pertanyaan where Pertanyaan=? ", quest.Pertanyaan)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer rows.Close()

	var result []Pertanyaan
	if rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.Pertanyaan, &each.Jawaban)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}

		result = append(result, each)
	}
	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	if result != nil {
		update, err := db.Exec("UPDATE Pertanyaan SET Jawaban = ? WHERE Pertanyaan = ?", quest.Jawaban, quest.Pertanyaan)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 2")
		}
		defer rows.Close()
		update.RowsAffected()
		return c.JSON(http.StatusOK, quest)
	}

	_, err = db.Exec("INSERT INTO pertanyaan (Pertanyaan,Jawaban) VALUES (?,?)", quest.Pertanyaan, quest.Jawaban)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusCreated, quest)
}

func addRespon(c echo.Context) error {
	var respon Respon
	err := c.Bind(&respon)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO respon (ID_histori,Jenis,Isi) VALUES (?,?,?)", respon.ID_histori, respon.Jenis, respon.Isi)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusCreated, respon)
}

func getChatFromId(c echo.Context) error {
	var histori HistoriReqId
	err := c.Bind(&histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}
	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}
	defer db.Close()
	rows, err := db.Query("select h.ID_histori, r.Jenis, r.Isi from histori as h, respon as r where h.ID_histori=r.ID_histori AND h.ID_histori=?", histori.ID_histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer rows.Close()
	var isi []Respon
	for rows.Next() {
		var respon Respon
		err := rows.Scan(&respon.ID_histori, &respon.Jenis, &respon.Isi)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}
		isi = append(isi, respon)
	}
	return c.JSON(http.StatusOK, isi)
}

func showHistori(c echo.Context) error {
	HistoriID := c.QueryParam("Id_histori")
	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}
	defer db.Close()
	fmt.Println(HistoriID)
	// HistoriID := "1"
	rows, err := db.Query("select h.ID_histori, r.Jenis, r.Isi from histori as h, respon as r where h.ID_histori=r.ID_histori AND h.ID_histori=?", HistoriID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer rows.Close()
	var isi []Respon
	for rows.Next() {
		var respon Respon
		err := rows.Scan(&respon.ID_histori, &respon.Jenis, &respon.Isi)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}
		isi = append(isi, respon)
	}
	fmt.Println(isi)
	nama, err := db.Query("select Nama from histori where ID_histori=?", HistoriID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 4")
	}
	defer nama.Close()
	var namaHistori string
	for nama.Next() {
		err := nama.Scan(&namaHistori)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 5")
		}
	}
	historiReq := &Histori{
		Nama: namaHistori,
		Isi:  isi,
	}
	return c.JSON(http.StatusOK, historiReq)
}

func addHistori(c echo.Context) error {
	var histori HistoriReq
	err := c.Bind(&histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO histori (Nama) VALUES (?)", histori.Nama)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusCreated, histori)
}

func deleteHistori(c echo.Context) error {
	var histori HistoriReq
	err := c.Bind(&histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	id, err := db.Query("select ID_histori from histori where Nama=?", histori.Nama)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer id.Close()
	var id_histori int
	for id.Next() {
		err := id.Scan(&id_histori)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}
	}
	_, err = db.Exec("DELETE FROM respon WHERE ID_histori=?", id_histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	_, err = db.Exec("DELETE FROM histori WHERE Nama=?", histori.Nama)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusCreated, histori)
}

// Knuth-Morris-Pratt algorithm
func KMP(text, pattern string) int {
	if len(pattern) == 0 {
		return 0
	}
	prefix := prefix(pattern)
	i := 0
	j := 0
	for i < len(text) {
		if text[i] == pattern[j] {
			i++
			j++
			if j == len(pattern) {
				return i - j
			}
		} else {
			if j == 0 {
				i++
			} else {
				j = prefix[j-1]
			}
		}
	}
	return -1
}

func prefix(pattern string) []int {
	prefix := make([]int, len(pattern))
	i := 1
	j := 0
	for i < len(pattern) {
		if pattern[i] == pattern[j] {
			prefix[i] = j + 1
			i++
			j++
		} else {
			if j == 0 {
				prefix[i] = 0
				i++
			} else {
				j = prefix[j-1]
			}
		}
	}
	return prefix
}

// Boyer-Moore algorithm
func BM(text, pattern string) int {
	if len(pattern) == 0 {
		return 0
	}
	last := last(pattern)
	i := len(pattern) - 1
	j := len(pattern) - 1
	for i < len(text) {
		if text[i] == pattern[j] {
			if j == 0 {
				return i
			} else {
				i--
				j--
			}
		} else {
			lo := last[int(text[i])]
			i = i + len(pattern) - min(j, 1+lo)
			j = len(pattern) - 1
		}
	}
	return -1
}

func last(pattern string) []int {
	last := make([]int, 256)
	for i := 0; i < 256; i++ {
		last[i] = -1
	}
	for i := 0; i < len(pattern); i++ {
		last[int(pattern[i])] = i
	}
	return last
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.GET("find", findAnswer)
	e.POST("quest", addQuestion)
	e.POST("respon", addRespon)
	e.GET("chat", showHistori)
	e.POST("histori", addHistori)
	e.GET("chatId", getChatFromId)
	e.DELETE("histori", deleteHistori)
	e.GET("allhistori", getAllHistori)
	e.Logger.Fatal(e.Start(":1234"))
}
