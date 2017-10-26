package main

import (
	"errors"
	"fmt"
	"github.com/kenmazsyma/soila/peer/api"
	cmn "github.com/kenmazsyma/soila/peer/common"
	"github.com/kenmazsyma/soila/peer/db"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

var appno int = 0

func main() {
	cmn.ReadEnv("soila.conf")
	if err := db.Init(); err != nil {
		fmt.Printf("DB Error:%s\n", err.Error())
		return
	}
	defer db.Term()
	if l, ok := startServer(); ok {
		defer l.Close()
	} else {
		return
	}
	procStdInput()
}

func startServer() (net.Listener, bool) {
	l, err := net.Listen("tcp", cmn.Env.SV_URL)
	if err != nil {
		log.Println("failed to run server:", err)
		return l, false
	}
	sv := rpc.NewServer()
	sv.Register(api.NewPerson())
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" {
			serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)
			err := sv.ServeRequest(serverCodec)
			if err != nil {
				fmt.Print(err)
				w.WriteHeader(500)
				return
			}
		} else {
			serveStatic(w, r)
		}
	}))
	return l, true
}

func procStdInput() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	exit_chan := make(chan int)
	go func() {
		s := <-signal_chan
		switch s {
		// kill -SIGHUP XXXX
		case syscall.SIGHUP:
			exit_chan <- 0
			return
		// kill -SIGINT XXXX or Ctrl+c
		case syscall.SIGINT:
			exit_chan <- 0
			return
		// kill -SIGTERM XXXX
		case syscall.SIGTERM:
			exit_chan <- 0
			return
		// kill -SIGQUIT XXXX
		case syscall.SIGQUIT:
			exit_chan <- 0
			return
		default:
			log.Println("Unknown signal.")
			return
		}
	}()
	code := <-exit_chan
	os.Exit(code)
}

var cont_type = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/x-javascript",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"svg":  "image/svg+xml",
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url == "/" {
		url = "/index.html"
	}
	exts := strings.Split(url, ".")
	ext := "html"
	prm := []string{}
	if len(exts) < 2 {
		url = url + ".html"
	} else {
		ext = exts[len(exts)-1]
		url = strings.Join(exts[:len(exts)-1], ".")
		prm = strings.Split(url, "_")
		url = prm[0] + "." + ext
	}
	data, err := readFile(url)
	if err != nil {
		w.WriteHeader(404)
		if ext == "html" || ext == "txt" {
			data, err = readFile("/error.html")
			if err == nil {
				w.Header().Set("Content-type", "text/html")
				w.Write(data)
				return
			}
			//	w.WriteHeader(404)
			w.Write([]byte("failed"))
		}
		return
	} else {
		w.Header().Set("Content-type", cont_type[ext])
		w.WriteHeader(200)
		w.Write(data)
	}
}

var static_data = map[string][]byte{}
var lock = sync.RWMutex{}

func readFile(path string) ([]byte, error) {
	lock.RLock()
	if cmn.Env.DEBUG == "0" {
		if buf, ok := static_data[path]; ok {
			lock.RUnlock()
			if len(buf) == 0 {
				fmt.Printf("failure to open the file : %s\n", path)
				return nil, errors.New("")
			}
			return buf, nil
		}
	}
	lock.RUnlock()
	fin, er := os.Open("html" + path)
	if er != nil {
		fmt.Printf("failure to open the file : %s\n", path)
		static_data[path] = []byte{}
		return nil, er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read the file : %s\n", path)
		static_data[path] = []byte{}
		return nil, er
	}
	lock.Lock()
	static_data[path] = buf
	defer lock.Unlock()
	return buf, nil
}
