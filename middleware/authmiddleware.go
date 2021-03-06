package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"ktsProject/models"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id := session.Value
		ses := models.Session{}
		query := "SELECT session::text, username::text, user_id::integer from session WHERE LOWER(session) = LOWER($1)"
		resultRows, _ := db.Query(query, id)
		defer resultRows.Close()
		for resultRows.Next() {
			err := resultRows.Scan(&ses.Session, &ses.Username, &ses.User_id)
			if err != nil {
				fmt.Println(err)
			}
		}
		if ses.Session == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		c := context.Background()
		ctx := context.WithValue(c, "UserId", ses.User_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
