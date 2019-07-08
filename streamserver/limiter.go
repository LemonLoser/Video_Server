package streamserver

import "log"

//bucket token 算法,实现流量限制
//recive one request, take token from bucket
//responsed,put the token back
type ConnectionLimit struct {
	MaxConnection int
	bucket        chan int
}

func NewConnectionLimit(cc int) *ConnectionLimit {
	return &ConnectionLimit{
		MaxConnection: cc,
		bucket:        make(chan int, cc),
	}
}

func (cl *ConnectionLimit) GetConn() bool {
	if len(cl.bucket) >= cl.MaxConnection {
		log.Printf("reached the rate limitation")
		return false
	}
	cl.bucket <- 1
	log.Printf("Successfully got connection")
	return true
}

func (cl *ConnectionLimit) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("new connection coming:%d", c)
}
