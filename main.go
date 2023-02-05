package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	maxUsers = 10
)

var (
	allUser    = make(map[net.Conn]string)
	connection = make(chan net.Conn)
)

func main() {
	var (
		port = flag.Int("p", 8989, "port")
		host = "localhost"
	)
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, *port))
	errorHandler(err)

	defer listener.Close()

	f, err := os.Create("data.txt")
	errorHandler(err)

	defer f.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil || len(allUser) == maxUsers {
				conn.Write([]byte("Access to this chat is closed\nPress any button to exit\n"))
				conn.Close()
				break
			}
			connection <- conn
		}
	}()

	for {
		go userLogin(<-connection, f)
	}
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkUserName(conn net.Conn, name string) (string, error) {
	flag := true
	for _, v := range name {
		if (v > 0 && v < 47) || (v > 'Z' && v < 'a') || v > 'z' {
			conn.Write([]byte("Username has invalid characters...\n[Please enter correct username]: "))
			flag = false
			break
		}
	}

	for _, value := range allUser {
		if value == name {
			conn.Write([]byte("Username is already taken...\n[Please enter other username]: "))
			flag = false
			break
		}
	}

	if len(name) == 0 {
		conn.Write([]byte("Username is empty...\n[ENTER YOUR NAME]: "))
		flag = false
	}

	if flag {
		return name, nil
	}

	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		conn.Write([]byte("Error Reader"))
		return "", err
	}

	name = name[:len(name)-1]
	return checkUserName(conn, name)
}

func userLogin(conn net.Conn, f *os.File) {
	content, err := ioutil.ReadFile("logo.txt")
	if err != nil {
		conn.Write([]byte("Welcome to TCP-Chat!"))
	} else {
		conn.Write(content)
	}

	conn.Write([]byte("[ENTER YOUR NAME]:"))
	name, err := bufio.NewReader(conn).ReadString('\n')

	if len(allUser) == maxUsers {
		conn.Write([]byte("Access to this chat is closed\nPress any button to exit\n"))
		conn.Close()
	}
	if err != nil {
		conn.Write([]byte("Error Reader"))
		return
	}
	name, err = checkUserName(conn, name[:len(name)-1])
	if err != nil {
		conn.Write([]byte("Error Reader"))
		return
	}

	conn.Write([]byte("\n"))
	allUser[conn] = name
	content, err = ioutil.ReadFile(f.Name())
	errorHandler(err)

	conn.Write(content)
	textMessage := "\r" + allUser[conn] + " has joined our chat...\n"

	for item := range allUser {
		msg := clear(userText(item)) + textMessage + userText(item)
		item.Write([]byte(msg))
	}

	if _, err := f.WriteString(textMessage); err != nil {
		log.Fatal()
	}

	go Chat(conn, f)
}

func userText(conn net.Conn) string {
	return fmt.Sprintf("\r[%s][%s]: ", time.Now().Format("01-02-2006 15:04:05"), allUser[conn])
}

func clear(a string) string {
	return "\r" + strings.Repeat(" ", len(a)) + "\r"
}

func Chat(conn net.Conn, f *os.File) {
	for {
		conn.Write([]byte(userText(conn)))
		message, err := bufio.NewReader(conn).ReadString('\n')
		if len(message) == 1 {
			continue
		}

		var msg string
		if err != nil {
			w := allUser[conn]
			delete(allUser, conn)
			textMessage := "\r" + w + " logged out of the chat...\n"
			for item := range allUser {
				msg := clear(userText(item)) + textMessage + userText(item)
				item.Write([]byte(msg))
			}
			if _, err := f.WriteString(textMessage); err != nil {
				log.Fatal()
			}
			return
		}

		textMessage := fmt.Sprintf("%s%s", userText(conn), string(message))

		for item := range allUser {
			if item == conn {
				continue
			}
			msg = clear(userText(item)) + textMessage

			item.Write([]byte(msg))
			time.Sleep(time.Second / 1000)
			item.Write([]byte(userText(item)))
		}

		if _, err := f.WriteString(textMessage); err != nil {
			log.Fatal()
		}
	}
}
