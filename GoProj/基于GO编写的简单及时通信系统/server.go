package main 

import(
	"fmt"
	"net"
	"sync"
	"io"
	"time"
)

type Server struct {
 	Ip string
	Port int

	OnlineMap map[string]*User
	mapLock sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip : ip,
		Port : port,
		OnlineMap: make(map[string]*User),
		Message : make(chan string), 
	}

	return server
}

func (this *Server) ListenMessager() {
	for{
		msg := <- this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" +user.Name + ":" + msg 
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("online sucesss")

	user := NewUser(conn, this)

	user.Online()

	// this.mapLock.Lock()
	// this.OnlineMap[user.Name] = user
	// this.mapLock.Unlock()

	// this.BroadCast(user, "onlined")
	isline := make(chan bool)

	go func (){
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0{
				//this.BroadCast(user, "downline")
				user.Offline()
				return
			}

			if err != nil && err != io.EOF{
				fmt.Println("conn Read err:", err)
				return 
			}

			msg := string(buf[:n-1])


			user.DoMessage(msg)
			isline <- true
			//this.BroadCast(user, msg)
		}
	}()

	for{
		select{
		case <-isline:

		case <-time.After(time.Second *5):
			user.SendMsg("you kik out")
			close(user.C)
			conn.Close()

			return
		}
	}

	//select {}
}

func (this *Server) Start (){
 	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",this.Ip, this.Port))
	if err != nil {
		fmt.Println("new.Listen err:", err)
		return 
	}
	
	defer listener.Close()

	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}

		go this.Handler(conn)
	}

}
