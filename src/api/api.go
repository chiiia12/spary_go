package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

type Result struct {
	Spa []Spa `json:"spa"`
}

type Spa struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func ShowSpaList(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@/" + os.Getenv("DB_NAME"))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close() // 関数がリターンする直前に呼び出される

	rows, err := db.Query("SELECT * FROM spa") //
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns() // カラム名を取得
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	//  rows.Scan は引数に `[]interface{}`が必要.

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	spaList := []Spa{}
	for rows.Next() {
		spa := Spa{}
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			if columns[i] == "name" {
				spa.Name = value
			} else if columns[i] == "address" {
				spa.Address = value
			} else if columns[i] == "id" {
				spa.Id = value
			}
		}
		spaList = append(spaList, spa)

	}
	data := Result{}
	data.Spa = spaList
	bytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s\n", string(bytes))
}

func Run() {
	http.HandleFunc("/api/spa/list", ShowSpaList)
	fmt.Printf("Server is running... localhost:8080")
	http.ListenAndServe(":8080", nil)

}
