package controllers

import (
	"encoding/json"
	"ktsProject/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (handle *Handle) TopCreateHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tx, err := handle.Db.Begin()
	if err != nil {
		log.Println(err)
	}
	id := r.Context().Value("UserId")
	Topic := models.Topic{}
	Topic.Creator_id = id.(int64)
	if r.FormValue("title") == "" || r.FormValue("body") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	Topic.Title = r.FormValue("title")
	Topic.Body = r.FormValue("body")
	Topic.Number_of_comments = 0
	Topic.Number_of_likes = 0
	Topic.Created = time.Now().String()
	insertTopicQuery := `insert into topics (title, body, number_of_comments, number_of_likes, creator_id, created ) values ($1, $2, $3, $4, $5, $6) returning id;`
	stmt, err := handle.Db.Prepare(insertTopicQuery)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(Topic.Title, Topic.Body, Topic.Number_of_comments, Topic.Number_of_likes, Topic.Creator_id, Topic.Created)
	err = row.Scan(&Topic.ID)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return
	}
	message, err := json.Marshal(Topic)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()
}

func (handle *Handle) TopListHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
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
	query := `select * from topics` + limit + offset
	resultRows, _ := handle.Db.Query(query)
	topics := make([]models.Topic, 0)
	defer resultRows.Close()
	for resultRows.Next() {
		Topic := models.Topic{}
		err := resultRows.Scan(&Topic.ID, &Topic.Title, &Topic.Body, &Topic.Number_of_comments, &Topic.Number_of_likes, &Topic.Creator_id, &Topic.Created)
		if err != nil {
			log.Println(err)
		}
		topics = append(topics, Topic)
	}
	message, err := json.Marshal(topics)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(message)
	if err != nil {
		log.Println(err)
	}
}

func (handle *Handle) TopLikeHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tx, err := handle.Db.Begin()
	if err != nil {
		log.Println(err)
	}
	TopicIdstr := r.FormValue("topic_id")
	check := handle.CheckTopic(TopicIdstr)
	if !check {
		handle.SendStatus(w, http.StatusBadRequest, "invalid input")
		return
	}
	TopicId, err := strconv.Atoi(TopicIdstr)
	if err != nil {
		log.Println(err)
	}
	user := r.Context().Value("UserId")
	query := "SELECT topic_id::integer, user_id::integer,created::timestamp  from topiclike WHERE topic_id = $1 and user_id = $2"
	like := models.TopicLike{}
	resultRows, err := handle.Db.Query(query, TopicId, user)
	if err != nil {
		log.Println(err)
	}
	for resultRows.Next() {
		err := resultRows.Scan(&like.Topic_id, &like.User_id, &like.Created)
		if err != nil {
			log.Println(err)
		}
	}
	defer resultRows.Close()
	interval := time.Minute
	if like.Topic_id != 0 && like.User_id != 0 {
		if err != nil {
			log.Println(err)
		}
		interval := interval.Minutes()
		delta := time.Now().UTC().Sub(like.Created.UTC()).Minutes()
		if delta <= interval {
			deleteQuery := `DELETE from topiclike WHERE topic_id = $1 and user_id = $2;`
			stmt, err := handle.Db.Prepare(deleteQuery)
			if err != nil {
				tx.Rollback()
				log.Println(err)
				return
			}
			defer stmt.Close()
			_, err = stmt.Exec(TopicId, user)
			if err != nil {
				log.Println(err)
			}
			updateTopicQuery := `update topics set number_of_likes=number_of_likes-1 WHERE id = $1;`
			stmt, err = handle.Db.Prepare(updateTopicQuery)
			if err != nil {
				tx.Rollback()
				log.Println(err)
				return
			}
			defer stmt.Close()
			_, err = stmt.Exec(TopicId)
			if err != nil {
				log.Println(err)
			}
			tx.Commit()
		} else {
			handle.SendStatus(w, http.StatusForbidden, "can not remove like")
			return
		}
	} else {
		insertLikeQuery := `insert into topiclike (topic_id, user_id, created ) values ($1, $2,$3 ) ON CONFLICT ON CONSTRAINT UNIQLIKE DO NOTHING;`
		stmt, err := handle.Db.Prepare(insertLikeQuery)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(TopicId, user, time.Now().UTC())
		if err != nil {
			log.Println(err)
		}
		updateTopicQuery := `update topics set number_of_likes=number_of_likes+1 WHERE id = $1;`
		stmt, err = handle.Db.Prepare(updateTopicQuery)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(TopicId)
		if err != nil {
			log.Println(err)
		}
		topicLike := models.TopicLike{int64(TopicId), user.(int64), time.Now().UTC()}
		w.WriteHeader(http.StatusOK)
		message, err := json.Marshal(topicLike)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(message)
		if err != nil {
			log.Println(err)
		}
		tx.Commit()
		return
	}
}
