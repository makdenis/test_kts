package routes

import (
	"github.com/gorilla/mux"
	"ktsProject/controllers"
	"ktsProject/middleware"
	"net/http"
)

func Router(handle *controllers.Handle) http.Handler {
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/auth.logout", handle.LogoutHandle)
	routerAuth.HandleFunc("/topic.create", handle.TopCreateHandle)
	routerAuth.HandleFunc("/topic.list", handle.TopListHandle)
	routerAuth.HandleFunc("/topic.like", handle.TopLikeHandle)
	routerAuth.HandleFunc("/comment.create", handle.CommentCreateHandle)
	routerAuth.HandleFunc("/comment.list", handle.CommentListHandle)
	authHandler := middleware.AuthMiddleware(routerAuth, handle.Db)

	router := mux.NewRouter()
	router.HandleFunc("/auth.login", handle.AuthHandle).Methods("POST")
	router.Handle("/topic.create", authHandler)
	router.Handle("/auth.logout", authHandler)
	router.Handle("/topic.list", authHandler)
	router.Handle("/topic.like", authHandler)
	router.Handle("/comment.create", authHandler)
	router.Handle("/comment.list", authHandler)
	return router
}
