package main

import(
	"fmt"
	"net"
	"flag"
	"os"
	"io"
)
type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp :serverIp,
		ServerPort: serverPort,
		flag:999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d",serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:",err)
	}

	client.conn = conn
	return client 
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("please input")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit"{ 
		if len(chatMsg) != 0{
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			break
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
	}

	chatMsg = ""
	fmt.Println("please input")
	fmt.Scanln(&chatMsg)
	
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	client.conn.Write([]byte(sendMsg))
}

func (client *Client) PrivateChat() {
	var remotename string
	var	chatMsg string
	
	client.SelectUser()
	fmt.Println("please user")
	fmt.Scanln(&remotename)

	for remotename != "exit" {
		fmt.Println("input msg")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0{
				sendMsg := "to|" + remotename + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				break
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
		}
	}

	chatMsg = ""
	fmt.Println("please user")
	fmt.Scanln(&remotename)
}

func (client *Client) menu() bool {

	var flag int

	fmt.Println("1.open line")
	fmt.Println("2.secret line")
	fmt.Println("3.change userName")
	fmt.Println("0.exit")

	fmt.Scanln(&flag)
	if flag >=0 && flag <=3 {
		client.flag = flag

		return true
	}else{
		return false
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)

	// for{
	// 	buf := make()
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (client *Client) UpdateName() bool {
	fmt.Println("input username:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name
	_, err := client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {
		case 1:
			// fmt.Println("open line")
			client.PublicChat()
		case 2:
			// fmt.Println("secert line")
			client.PrivateChat()
		case 3:
			// fmt.Println("userName change")
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&cd, "ip", "127.0.0.1", "set server addr")
	flag.IntVar(&serverPort, "port", 8888, "set server port")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp,serverPort)
	if client == nil{
		fmt.Println("connect failed")
	}

	go client.DealResponse()

	fmt.Println("connect sucess........")

	client.Run()
	select{}
}
