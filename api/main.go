package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"videoProject/api/session"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//api check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}
func RegisterHandlers() *httprouter.Router {
	log.Printf("preparing to post request\n")
	router := httprouter.New()
	router.POST("/user", CreateUser) //CreateUser作为一个function传入
	router.POST("/user/:user_name", Login)
	router.GET("/user/:user_name", GetUserInfo)
	router.POST("/user/:username/videos", AddNewVideo)
	router.GET("/user/:username/videos", ListAllVideos)
	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
	router.POST("/videos/:vid-id/comments", PostComment)
	router.GET("/videos/:vid-id/comments", ShowComments)
	return router
}

func Prepare() {
	session.LoadSessionsFromDB()
}
func main() {
	Prepare()
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh) //lisen->RegisterHandlers->Handlers
}

//Handler->validation{1.request 2.user}校验->business logic->response
//main->middleware->defs(message,err)->handlers->dbops->response
