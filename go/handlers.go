package main

import (
	"fmt"
	"net"
	"strings"
)

func pingHandler(conn net.Conn, args []string) {
	code := args[1]

	sendCommand(conn, "PONG", code)
	fmt.Println("Pinged and Ponged, code is " + code)
}

func errorHandler(conn net.Conn, args []string) {
	err := args[1]

	fmt.Println("~~Error occured: " + err)
	//reloadConnection()
}

func noticeHandler(conn net.Conn, args []string) {
	//text := args[1]
	//fmt.Println("~~~Notice Handler")
}

func privHandler(conn net.Conn, args []string) {
	nick := args[1]
	//user := args[2]
	//host := args[3]
	dest := args[4]
	mesg := args[5]

	//fmt.Println("~~~Priv Handler")

	responseDest := dest

	if dest == nickname {
		responseDest = nick
	}

	switch {
	case strings.HasPrefix(mesg, "!echo"):
		mesg = strings.Replace(mesg, "!echo", "", 1)
		mesg = strings.TrimSpace(mesg)
		sendPrivMsg(conn, responseDest, mesg)
	case strings.HasPrefix(mesg, "!stop"):
		conn.Close()
	case checkWeatherCommand(mesg):
		printWeatherData(conn, responseDest, mesg)
	}

}

func comdHandler(conn net.Conn, args []string) {
	number := args[1]
	//who := args[2]
	//text := args[3]

	fmt.Println("~~~COMD Handler number is " + number)

	switch number {
	case "443":
		fmt.Println("~~Nick in use trying again~~")
	case "001":
		joinChannels(conn)
		fmt.Println("~~~Join Handler")
	}
}

func joinHandler(conn net.Conn, args []string) {
	//nick := args[1]
	//user := args[2]
	//host := args[3]
	//channel := args[4]
	fmt.Println("~~~Join Handler")
}

func quitHandler(conn net.Conn, args []string) {
	//user := args[1]
	//reason := args[2]
	fmt.Println("~~~Quit Handler")
}
