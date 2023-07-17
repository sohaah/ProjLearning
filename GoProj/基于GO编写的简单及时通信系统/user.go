package main

import(
	"net"
	"io"
	"fmt"
	"strings"
)

type User struct{
	Name string
	Addr string
	C chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr, 
		C:	  make(chan string),
		conn: conn,
		server : server,
	}

	go user.ListenMessager()

	return user
}

func (this *User) Online() {
		//fmt.Println("online sucesss")

		this.server.mapLock.Lock()
		this.server.OnlineMap[this.Name] = this
		this.server.mapLock.Unlock()
	
		this.server.BroadCast(this, "onlined")
	
		go func (){
			buf := make([]byte, 4096)
			for {
				n, err := this.conn.Read(buf)
				if n == 0{
					this.server.BroadCast(this, "downline")
					return
				}
	
				if err != nil && err != io.EOF{
					fmt.Println("conn Read err:", err)
					return 
				}
	
				msg := string(buf[:n-1])
				fmt.Println(msg)
				//this.server.BroadCast(this, msg)
				this.DoMessage(msg)
			}
		}()
	
		select {}
}

func (this *User) Offline() {
	
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "offline")

}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
} 

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name +":" + "online1111" + "\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	}else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("name already in use")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			
			this.Name = newName
			this.SendMsg("name already change")
		}
	}else if len(msg) > 4 && msg[:3] == "to|"{
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("fmt wrong")
		}
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("user none")
			return
		}
		
		content := strings.Split(msg, "|")[2]
		remoteUser.SendMsg(this.Name + "said" + content)
	

	}else{
		fmt.Println("8888888888888888888")
		this.server.BroadCast(this, msg)
	}
}

func (this *User) ListenMessager() {
	for{
		msg := <- this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
