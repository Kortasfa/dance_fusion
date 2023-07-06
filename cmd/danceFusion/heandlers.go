package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type buttonRequest struct {
	isBtnClicked bool `json:"button"`
}

func gameField(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/gameField.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func gamePhone(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/gamePhone.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func checkButton(w http.ResponseWriter, r *http.Request) {
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	var req buttonRequest

	err = json.Unmarshal(reqData, &req)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	w.WriteHeader(200)

}

/*func start(w http.ResponseWriter, r *http.Request) {
	emp := &buttonRequest{isBtnClicked: true}
	e, err := json.Marshal(emp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

	w.WriteHeader(200)
}*/
