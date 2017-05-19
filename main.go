package main

import (
	"flag"
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:1996", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var err error

var m map[string][]*websocket.Conn

func echo(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//check
	c, ok := m[r.URL.Path]

	if ok {
		c = append(c, conn)
		m[r.URL.Path] = c
		fmt.Println("追加指针")
	} else {
		var c []*websocket.Conn
		c = append(c, conn)
		m[r.URL.Path] = c
		fmt.Println("创建第一个指针")
	}

	log.Println(r.URL.Path)

	defer conn.Close()
	for {

		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		//分发数据
		for _, s := range m[r.URL.Path] {
			if s == conn {
				fmt.Println("跳过当前建立链接的指针")
				continue
			}
			fmt.Println("正在分发数据")
			err = s.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func main() {
	m = make(map[string][]*websocket.Conn)
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
