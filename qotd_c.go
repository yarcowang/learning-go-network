/**
 * @file qotd_c.go
 * it is a qotd client, nothing more... :)
 * cmd: qotd_c [number of quotes you want]
 *
 * @author Yarco <yarco.wang@gmail.com>
 * @since 2016-09-10
 * @copyright BSD
 */
//  vim: set tabstop=2 shiftwidth=2 softtabstop=2 noexpandtab ai si: 
package main

import (
	"net"
	"fmt"
	"os"
	"strconv"
)

const (
	// nicksosinski.com server from wiki
	QOTD_QUOTES_WANT_DEFAULT = 3
	QOTD_PROTO = "tcp"
	QOTD_SERVER = "nicksosinski.com"
	QOTD_PORT = "17"
)

func MustBe(obj interface{}, err error) interface{} {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %s\n", err.Error())
		os.Exit(1)
	}

	return obj
}

func ShowMeTheQuote() string {
	conn := MustBe(net.Dial(QOTD_PROTO, QOTD_SERVER + ":" + QOTD_PORT)).(net.Conn)
	defer func() {
		// fmt.Println("Close connection...Bye bye.")
		conn.Close()
	}()

	msg := make([]byte, 1024)
	n := MustBe(conn.Read(msg)).(int)
	return string(msg[:n])
}

func main() {
	var n int

	if len(os.Args) != 1 {
		n = MustBe(strconv.Atoi(os.Args[1])).(int)
	} else {
		n = QOTD_QUOTES_WANT_DEFAULT
	}

	ch := make(chan string, n)
	for i := 0; i < n; i++ { // is it flooding?
		go func() {
			ch <- ShowMeTheQuote()
		}()
	}

	for i := 0; i < n; i++ {
		fmt.Println(<- ch)
	}
}

