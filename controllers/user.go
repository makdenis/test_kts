package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"ktsProject/models"
	"log"
	"net/http"
	"time"
)

func (handle *Handle) AuthHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tx, err := handle.Db.Begin()
	if err != nil {
		log.Println(err)
	}
	user := models.User{}
	query := "SELECT id::integer, username::text, password::text,email::text, first_name::text, last_name::text, date_joined::text  from users WHERE LOWER(username) = LOWER($1)"
	resultRows, _ := handle.Db.Query(query, r.FormValue("username"))
	defer resultRows.Close()
	for resultRows.Next() {
		err := resultRows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.First_name, &user.Last_name, &user.Date_joined)
		if err != nil {
			log.Println(err)
		}
	}
	if user.Username == "" {
		handle.SendStatus(w, http.StatusBadRequest, "invalid input")
		return
	}
	if user.Password != r.FormValue("password") {
		handle.SendStatus(w, http.StatusUnauthorized, "wrong nick or pass")
		return
	}
	session := uuid.New().String()
	insertUserQuery := `insert into session (user_id, username, session) values ($1, $2, $3);`
	stmt, err := handle.Db.Prepare(insertUserQuery)
	if err != nil {
		tx.Rollback()
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID, user.Username, session)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    session,
		Expires:  time.Now().Add(60 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	message, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()
	w.Write(message)
}

func (handle *Handle) LogoutHandle(w http.ResponseWriter, r *http.Request) {
	tx, err := handle.Db.Begin()
	if err != nil {
		log.Println(err)
	}
	session, err := r.Cookie("sessionid")
	if err != nil {
		log.Println(err)
	}
	id := session.Value
	deleteQuery := `DELETE from session WHERE LOWER(session) = LOWER($1);`
	stmt, err := handle.Db.Prepare(deleteQuery)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(deleteQuery, id)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:     "sessionid",
		Value:    id,
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	tx.Commit()
}
