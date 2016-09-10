/**
 * @file rfc.go
 * rfc reader
 * a tool to fetch rfc text from
 *	https://tools.ietf.org/rfc/rfc[\d+].txt
 * cache it to your home directory and output
 *
 * @author Yarco <yarco.wang@gmail.com>
 * @since 2016-09-10
 * @copyright BSD
 */
//  vim: set tabstop=2 shiftwidth=2 softtabstop=2 noexpandtab ai si: 
package main

import (
	"io"
	"os"
	"os/user"
	"net/http"
	"path"
	"fmt"
)

const (
	IETF_RFC_URL = "https://tools.ietf.org/rfc/rfc%s.txt"
)

func MustBe(obj interface{}, err error) interface{} {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %s\n", err.Error())
		os.Exit(1)
	}

	return obj
}

func myRfcHome() string {
	u := MustBe(user.Current()).(*user.User)
	return path.Join(u.HomeDir, ".rfcs")
}

func cacheRfc(rfc_remote string, rfc_local string) {
	// fetch rfc
	resp := MustBe(http.Get(rfc_remote)).(*http.Response)
	toFile := MustBe(os.Create(rfc_local)).(*os.File)
	defer func() {
		resp.Body.Close()
		toFile.Close()
	}()

	MustBe(io.Copy(toFile, resp.Body))
}

func main() {
	if (len(os.Args) != 2) {
		fmt.Printf("%s <rfc no.>\n  Example: %s 1024\n", os.Args[0], os.Args[0])
		return
	}

	// create the directory
	home := myRfcHome()
	os.MkdirAll(home, 0777)

	// rfc file
	rfc := path.Join(home, os.Args[1])
	if _, err := os.Stat(rfc); os.IsNotExist(err) { // fetch and create the file
		cacheRfc(fmt.Sprintf(IETF_RFC_URL, os.Args[1]), rfc)
	}

	f := MustBe(os.Open(rfc)).(*os.File)
	defer f.Close()

	io.Copy(os.Stdout, f)
}



