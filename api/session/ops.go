package session

import (
	"fmt"
	"sync"
	"time"
	"videoProject/api/dbops"
	"videoProject/api/defs"
	"videoProject/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionsFromDB() *sync.Map {
	r, err := dbops.RetrieveAllSession()
	if err != nil {
		return nil
	}
	r.Range(func(key, value interface{}) bool {
		ss := value.(*defs.SimpleSession)
		sessionMap.Store(key, ss)
		return true
	})
	return sessionMap
}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := time.Now().UnixNano() / 1000000
	ttl := ct + 30*60*1000 //30min
	ss := &defs.SimpleSession{Username: un, TTL: ttl}
	sessionMap.Store(id, ss)
	err := dbops.InsertSession(id, ttl, un)
	if err != nil {
		return fmt.Sprintf("error of generateNewSessionId:%s", err)
	}
	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := time.Now().UnixNano() / 100000
		if ss.(*defs.SimpleSession).TTL < ct {
			//delete expired session
			DeleteExpiredSession(sid)
			return "", true
		}
		return ss.(defs.SimpleSession).Username, false
	}
	return "", true
}

func DeleteExpiredSession(sid string) error {
	sessionMap.Delete(sid)
	err := dbops.DeleteSession(sid)
	if err != nil {
		return err
	}
	return nil
}
