package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
)

const hostname = ""
const port = ""
const password = ""

var channels = [...]string{
	"#nickbot",
	"#squad"}

const nickname = "SystemD"

//connect ...
//Creates the inital connection to the irc server
func connect() net.Conn {
	//connection, err := net.Dial("unix", "/var/www/bot/irc/sock")
	connection, err := net.Dial("tcp", "localhost:8765")

	//	sendCommand(connection, "PASS", password)
	//	sendCommand(connection, "USER", fmt.Sprintf("%s 8 * :%s\r\n", nickname, nickname))
	//	sendCommand(connection, "NICK", nickname)
	//	fmt.Println("Registration sent")

	if err != nil {
		panic(err)
	}

	joinChannels(connection)

	return connection
}

//output ...
//Output anything the server sends us and process it
func listen(conn net.Conn) {

    fmt.Println("Starting Listener")
	scanner := bufio.NewScanner(conn)

	rePing := regexp.MustCompile(`PING :(.+)$`)
	reError := regexp.MustCompile(`ERROR :(.+)$`)
	reNotice := regexp.MustCompile(`:[^ ]+ NOTICE [^ ]+ :(.+)$`)
	reComd := regexp.MustCompile(`:[^ ]+ (\d\d\d) ([^:]+)(?: :(.+))?$`)
	rePriv := regexp.MustCompile(`:([^!]+)!([^@]+)@([^ ]+) PRIVMSG ([^ ]+) :(.+)$`)
    reJoin := regexp.MustCompile(`:([^!]+)!([^@]+)@([^ ]+) JOIN :([^ ]+)$`)
	reQuit := regexp.MustCompile(`:([^!]+)[^ ]+ QUIT :(.+)$`)

	find := func(reg *regexp.Regexp, msg string) []string {
		return reg.FindStringSubmatch(msg)
	}

    fmt.Println("Starting scanner")
	for scanner.Scan() {
		msg := scanner.Text()

		fmt.Println(string(msg) + "~")

		switch {
		case rePing.MatchString(msg):
			go pingHandler(conn, find(rePing, msg))

		case reError.MatchString(msg):
			go errorHandler(conn, find(reError, msg))

		case reNotice.MatchString(msg):
			go noticeHandler(conn, find(reNotice, msg))

		case reComd.MatchString(msg):
			go comdHandler(conn, find(reComd, msg))

		case rePriv.MatchString(msg):
			go privHandler(conn, find(rePriv, msg))

		case reJoin.MatchString(msg):
			go joinHandler(conn, find(reJoin, msg))

		case reQuit.MatchString(msg):
			go quitHandler(conn, find(reQuit, msg))

		}
	}
}

func joinChannels(conn net.Conn) {
	fmt.Printf("Joining channels %s", channels)
	for _, c := range channels {
		sendCommand(conn, "JOIN", c)
		fmt.Println("Joining channel " + c)
	}
}

//printIRC ...
//Helper method to send the command and text to the irc server
func sendCommand(conn net.Conn, command string, text string) {
	fmt.Fprintf(conn, "%s %s\r\n", command, text)
    fmt.Printf("%s %s\n", command, text)
}

func sendPrivMsg(conn net.Conn, channel string, text string) {
	sendCommand(conn, "PRIVMSG", fmt.Sprintf("%s :%s\r\n", channel, text))
}

func main() {
	conn := connect()

	defer conn.Close()

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		conn.Close()

		os.Exit(0)
	}()

	listen(conn)
}
