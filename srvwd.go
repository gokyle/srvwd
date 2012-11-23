package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

const VERSION = "2.0.0"

var (
	srv_addr   string
	srv_bin    string
	srv_chroot bool
	srv_wd     string
	srv_port   int
	srv_usr    int
	srv_ver    bool
)

func init() {
	srv_bin = filepath.Base(os.Args[0])
	config()
}

func main() {
	srv_addr := fmt.Sprintf(":%d", srv_port)
	fmt.Printf("serving %s on %s\n", srv_wd, srv_addr)
	http.Handle("/", http.FileServer(http.Dir(srv_wd)))
	log.Fatal(http.ListenAndServe(srv_addr, nil))
}

func config() {
	empty := func(s string) bool {
		return len(s) == 0
	}

	fPort := flag.Int("p", 8080, "port to listen on")
	fChroot := flag.Bool("r", false, "chroot to the working directory")
	fUser := flag.String("u", "", "user to run as")
	fVersion := flag.Bool("v", false, "print version information")
	flag.Parse()

	if *fVersion {
		version()
	}

	fmt.Printf("flag.args: %+v\n", flag.Args())
	if flag.NArg() > 0 {
		srv_wd = flag.Arg(0)
	}

	if !empty(*fUser) {
		setuid(*fUser)
	}

	if *fChroot {
		chroot(srv_wd)
	}

	srv_port = *fPort
}

func fatal(err error) {
	fmt.Printf("[!] %s: %s\n", srv_bin, err.Error())
}

func checkFatal(err error) {
	if err != nil {
		fatal(err)
	}
}

func setuid(username string) {
	usr, err := user.Lookup(username)
	checkFatal(err)
	uid, err := strconv.Atoi(usr.Uid)
	checkFatal(err)
	err = syscall.Setreuid(uid, uid)
	checkFatal(err)
}

func chroot(path string) {
	err := syscall.Chroot(path)
	checkFatal(err)
	srv_wd = "/"
}

func version() {
	fmt.Printf("%s version %s\n", srv_bin, VERSION)
	os.Exit(0)
}
