package main

import (
	"database/sql"
	"fmt"
	"net/http"
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
	PertayaanHistori struct {
		Pertanyaan string `json:"pertanyaan"`
		ID_histori int    `json:"id_histori"`
	}
	HistoriReq struct {
		Nama string `json:"nama"`
	}

	Histori struct {
		Nama string   `json:"nama"`
		Isi  []Respon `json:"isi"`
	}
	UpdateHistori struct {
		NewName    string `json:"new_name"`
		ID_histori int    `json:"ID_histori"`
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
	db, err := sql.Open("mysql", "root:qwerty@tcp(127.0.0.1:3306)"+"/tubes3")
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
	var questHistori PertayaanHistori
	err := c.Bind(&questHistori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	var quest PertanyaanReq
	quest.Pertanyaan = questHistori.Pertanyaan

	//add respon question
	var statusRespon string
	var respon Respon
	respon.ID_histori = questHistori.ID_histori
	respon.Jenis = "input"
	respon.Isi = quest.Pertanyaan
	statusRespon = addResponReq(respon)
	if statusRespon == "success" {
		fmt.Println("Berhasil add respon")
	} else {
		fmt.Println("Gagal add respon")
	}

	quest.Pertanyaan = strings.ToLower(quest.Pertanyaan)

	if dateCheck(quest.Pertanyaan) {
		var statusRespon string
		var respon Respon
		respon.ID_histori = questHistori.ID_histori
		respon.Jenis = "output"
		respon.Isi = "Hari dari tanggal " + parsingDate(quest.Pertanyaan) + " adalah " + getDay(parsingDate(quest.Pertanyaan))
		statusRespon = addResponReq(respon)
		if statusRespon == "success" {
			fmt.Println("Berhasil add respon")
		} else {
			fmt.Println("Gagal add respon")
		}

		return c.JSON(http.StatusOK, "Hari dari tanggal "+parsingDate(quest.Pertanyaan)+" adalah "+getDay(parsingDate(quest.Pertanyaan)))
	} else if updateQuestionCheck(quest.Pertanyaan) {
		var question Pertanyaan
		question.Jawaban = parsingUpdateQuestion(quest.Pertanyaan)[1]
		question.Pertanyaan = parsingUpdateQuestion(quest.Pertanyaan)[0]
		var status string
		status = addQuestionReq(question)

		var statusRespon string
		var respon Respon
		respon.ID_histori = questHistori.ID_histori
		respon.Jenis = "output"

		if status == "success" {
			if questionCheck(question.Pertanyaan) {
				respon.Isi = "Pertanyaan " + question.Pertanyaan + " sudah ada! Jawaban diupdate ke " + question.Jawaban
			} else {
				respon.Isi = "Pertanyaan " + question.Pertanyaan + "telah ditambahkan"
			}
			statusRespon = addResponReq(respon)

			if statusRespon == "success" {
				fmt.Println("Berhasil add respon")
			} else {
				fmt.Println("Gagal add respon")
			}

			return c.JSON(http.StatusOK, "Berhasil update pertanyaan")
		} else {
			return c.JSON(http.StatusOK, "Gagal update pertanyaan")
		}
	} else if deleteQuestionCheck(quest.Pertanyaan) {
		var question string
		question = parsingDeleteQuestion(quest.Pertanyaan)
		var status string
		status = deleteQuestionReq(question)
		if questionCheck(question) {
			var statusRespon string
			var respon Respon
			respon.ID_histori = questHistori.ID_histori
			respon.Jenis = "output"
			respon.Isi = "Pertanyaan " + question + " telah dihapus"
			statusRespon = addResponReq(respon)
			if statusRespon == "success" {
				fmt.Println("Berhasil add respon")
			} else {
				fmt.Println("Gagal add respon")
			}
		} else {
			var statusRespon string
			var respon Respon
			respon.ID_histori = questHistori.ID_histori
			respon.Jenis = "output"
			respon.Isi = "Tidak ada pertanyaan " + question + " dalam database"
			statusRespon = addResponReq(respon)
			if statusRespon == "success" {
				fmt.Println("Berhasil add respon")
			} else {
				fmt.Println("Gagal add respon")
			}
		}
		if status == "success" {
			return c.JSON(http.StatusOK, "Berhasil hapus pertanyaan")
		} else {
			return c.JSON(http.StatusOK, "Gagal hapus pertanyaan")
		}
	} else {
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

}

func questionCheck(question string) bool {
	db, err := connect()
	if err != nil {
		fmt.Println("error")
	}
	defer db.Close()

	rows, err := db.Query("select Pertanyaan,Jawaban from pertanyaan where Pertanyaan=? ", question)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var result []Pertanyaan
	fmt.Println(rows)
	if rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.Pertanyaan, &each.Jawaban)
		if err != nil {
			fmt.Println(err.Error())
		}
		result = append(result, each)
	}
	if len(result) > 0 {
		return true
	} else {
		return false
	}
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
		fmt.Println(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer rows.Close()

	var result []Pertanyaan
	fmt.Println(rows)
	if rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.Pertanyaan, &each.Jawaban)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
		}

		result = append(result, each)
	}
	fmt.Println(result)
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

func deleteQuestionReq(question string) string {
	db, err := connect()
	if err != nil {
		return "err"
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM pertanyaan WHERE Pertanyaan = ?", question)
	if err != nil {
		return "err"
	}

	return "success"
}

func addQuestionReq(quest Pertanyaan) string {
	db, err := connect()
	if err != nil {
		return "err"
	}
	defer db.Close()

	rows, err := db.Query("select Pertanyaan,Jawaban from pertanyaan where Pertanyaan=? ", quest.Pertanyaan)
	if err != nil {
		fmt.Println(err.Error())
		return "err"
	}
	defer rows.Close()

	var result []Pertanyaan
	fmt.Println(rows)
	if rows.Next() {
		var each = Pertanyaan{}
		var err = rows.Scan(&each.Pertanyaan, &each.Jawaban)

		if err != nil {
			fmt.Println(err.Error())
			return "err"
		}

		result = append(result, each)
	}
	fmt.Println(result)
	if err = rows.Err(); err != nil {
		return "err"
	}

	if result != nil {
		update, err := db.Exec("UPDATE Pertanyaan SET Jawaban = ? WHERE Pertanyaan = ?", quest.Jawaban, quest.Pertanyaan)
		if err != nil {
			return "err"
		}
		defer rows.Close()
		update.RowsAffected()
		return "success"
	}

	_, err = db.Exec("INSERT INTO pertanyaan (Pertanyaan,Jawaban) VALUES (?,?)", quest.Pertanyaan, quest.Jawaban)
	if err != nil {
		return "err"
	}

	return "success"
}
func addResponReq(respon Respon) string {

	db, err := connect()
	if err != nil {
		return "err"
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO respon (ID_histori,Jenis,Isi) VALUES (?,?,?)", respon.ID_histori, respon.Jenis, respon.Isi)
	if err != nil {
		return "err"
	}

	return "success"
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
	nama, err := db.Query("select Nama from histori where ID_histori=?", HistoriID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer nama.Close()
	var namaHistori string
	for nama.Next() {
		err := nama.Scan(&namaHistori)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, "error 3")
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
	id_histori := c.QueryParam("Id_histori")
	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}
	_, err = db.Exec("DELETE FROM respon WHERE ID_histori=?", id_histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	_, err = db.Exec("DELETE FROM histori WHERE ID_histori=?", id_histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 3")
	}

	return c.JSON(http.StatusCreated, "successs")
}

func updateHistoriName(c echo.Context) error {
	var historiChange UpdateHistori
	err := c.Bind(&historiChange)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 1")
	}

	db, err := connect()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	defer db.Close()

	update, err := db.Exec("UPDATE Histori SET Nama = ? WHERE ID_histori = ?", historiChange.NewName, historiChange.ID_histori)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "error 2")
	}
	update.RowsAffected()
	return c.JSON(http.StatusOK, "success")
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
	e.GET("histori", showHistori)
	e.POST("histori", addHistori)
	e.GET("chat", getChatFromId)
	e.DELETE("histori", deleteHistori)
	e.GET("allhistori", getAllHistori)
	e.PUT("histori", updateHistoriName)
	e.Logger.Fatal(e.Start(":1234"))
}
