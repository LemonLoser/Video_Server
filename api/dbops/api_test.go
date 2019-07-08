package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempvid string

//intit(dblogin,truncate tables)-->run test-->clear data(truncate tables)
func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate vedio_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}
func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Delete", testDeleteUser)
	t.Run("Reget", testReGetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("acivc", "123")
	if err != nil {
		t.Errorf("error of addUser %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("acivc")
	if pwd != "123" || err != nil {
		t.Errorf("error of getUser %v", err)
	}
}

func testDeleteUser(t *testing.T) {
	err = DeleteUser("acivc", "123")
	if err != nil {
		t.Errorf("error of deleteUser %v", err)
	}
}

//检测是否删除成功
func testReGetUser(t *testing.T) {
	pwd, err := GetUserCredential("acivc")
	if err != nil {
		t.Errorf("error of RegetUser %v", err)
	}
	if pwd != "" {
		t.Errorf("deleting user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("prepareUser", testAddUser)
	t.Run("addvideo", testAddNewVideo)
	t.Run("getvideo", testGetVideoInfo)
	t.Run("deletevideo", testDeleteVideoInfo)
	t.Run("regetvideo", testRegetVideoInfo)
}
func testAddNewVideo(t *testing.T) {
	v, err := AddNewVideo("1", "my_info")
	if err != nil {
		t.Errorf("error of addNewVedio %v", err)
	}
	tempvid = v.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("error of getVedioInfo %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("error of detele")
	}
}

func testRegetVideoInfo(t *testing.T) {
	v, err := GetVideoInfo(tempvid)
	if err != nil || v != nil {
		t.Errorf("error of regetVedioInfo %v", err)
	}
}

func TestComment(t *testing.T) {
	clearTables()
	t.Run("Adduser", testAddUser)
	t.Run("AddComments", testAddNewComments)
	t.Run("ListComments", testListComments)
}

func testAddNewComments(t *testing.T) {
	vid := "123"
	aid := "1"
	content := "video"
	err := AddNewComments(vid, aid, content)
	if err != nil {
		t.Errorf("error of addNewComment %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 123456780
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("error of listcomments %v", err)
	}
	for index, ele := range res {
		fmt.Printf("comments:%d,%v \n", index, ele)
	}

}
