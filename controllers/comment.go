package controllers

import (
	"encoding/json"
	"ktsProject/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (handle *Handle) CommentCreateHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := r.FormValue("body")
	topic_idstr := r.FormValue("topic_id")
	check := handle.CheckTopic(topic_idstr)
	if !check {
		handle.SendStatus(w, http.StatusBadRequest, "invalid input")
		return
	}
	TopicId, err := strconv.Atoi(topic_idstr)
	if err != nil {
		log.Println(err)
	}
	if body == "" {
		handle.SendStatus(w, http.StatusBadRequest, "invalid input")
		return
	}
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session := cookie.Value
	comment := models.Comment{}
	id := handle.GetUserId(session)
	comment.Creator_id = id
	comment.Topic_id = int64(TopicId)
	comment.Body = body
	comment.Created = time.Now().UTC()
	insertCommentQuery := `insert into comments (body, creator_id, created, topic_id ) values ($1, $2, $3, $4) returning id;`
	row := handle.Db.QueryRow(insertCommentQuery, comment.Body, comment.Creator_id, comment.Created, comment.Topic_id)
	err = row.Scan(&comment.ID)
	if err != nil {
		log.Println(err)
	}
	updateTopicQuery := `update topics set number_of_comments=number_of_comments+1 WHERE id = $1;`
	_, err = handle.Db.Exec(updateTopicQuery, comment.Topic_id)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	message, err := json.Marshal(comment)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
	}
}

func (handle *Handle) CommentListHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	topic_id := r.FormValue("topic_id")
	if limit == "" {
		limit = " limit 100"
	} else {
		if !handle.CheckInput(limit) {
			handle.SendStatus(w, http.StatusBadRequest, "invalid input")
			return
		}
		limit = " limit " + limit
	}
	if offset != "" {
		if !handle.CheckInput(offset) {
			handle.SendStatus(w, http.StatusBadRequest, "invalid input")
			return
		}
		offset = " offset " + offset
	}
	check := handle.CheckTopic(topic_id)
	if !check {
		handle.SendStatus(w, http.StatusBadRequest, "invalid input")
		return
	}
	query := `select id, body, creator_id, created from comments where topic_id= $1` + limit + offset
	resultRows, _ := handle.Db.Query(query, topic_id)
	comments := make([]models.Comment, 0)
	defer resultRows.Close()
	for resultRows.Next() {
		comment := models.Comment{}
		err := resultRows.Scan(&comment.ID, &comment.Body, &comment.Creator_id, &comment.Created)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, comment)
	}
	message, err := json.Marshal(comments)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
	}
}
