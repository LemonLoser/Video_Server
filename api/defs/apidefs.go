package defs

//requests
type UserCredential struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type NewComment struct {
	AuthorId string `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

//response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserSession struct {
	UserName  string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Id int `json:"id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}

//Data model
type User struct {
	Id        int    `json:"id"`
	LoginName string `json:"login_name"`
	Pwd       string `json:"pwd"`
}

type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     string `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

type Comment struct {
	Id      string `json:"id"`
	VideoId string `json:"video_id"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type SimpleSession struct {
	Username string //login name
	TTL      int64
}
