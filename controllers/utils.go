package controllers

import (
	"encoding/json"
	"fmt"
	"ktsProject/models"
	"log"
	"net/http"
	"strconv"
)


func (handle *Handle) CheckTopic(id string) bool {
	ok := 0
	query := "SELECT id::integer  from topics WHERE id = $1"
	resultRows, _ := handle.Db.Query(query, id)
	defer resultRows.Close()
	for resultRows.Next() {
		err := resultRows.Scan(&ok)
		if err != nil {
			log.Println(err)
		}
	}
	if ok != 0 {
		return true
	} else {
		return false
	}
}

func (handle *Handle) CheckInput(digit string) bool {
	_, err := strconv.Atoi(digit)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (handle *Handle) SendStatus(w http.ResponseWriter, status int, mes string) {
	w.WriteHeader(status)
	message, err := json.Marshal(models.Err{mes})
	if err != nil {
		fmt.Println(err)
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
	}
}
