package main

import (
	"log"
	"net"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var stdPong = int64(0)

func HandleConn(c net.Conn) {
	defer c.Close()

	// handle incoming data
	buffer := make([]byte, 1024)
	numBytes, err := c.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("received", numBytes, "bytes:", string(buffer))

	// handle reply
	msg := string(buffer[:numBytes]) + " -> " + time.Now().Format("15:04:05.000000")
	_, err = c.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func Ping(proto, addr string) string {

	//log.Println(addr)
	c, err := net.Dial(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	msg := []byte(time.Now().Format("15:04:05.000000"))
	_, err = c.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	_, err = c.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func RenderBar(diff int64) string {
	var b int

	b = int(float64(diff) / float64(stdPong) * 10)
	s := "["
	for i := 0; i < 30; i++ {
		if i < b {
			s += "|"
		} else {
			s += "."
		}
	}
	s += "]"
	return s

}

func main() {

	var (
		listenAddress = kingpin.Flag("p", "Address to listen").Default(":8356").String()
		appType       = kingpin.Flag("t", "Ping or pong ? ").Default("pong").String()
		dstAddress    = kingpin.Flag("d", "Destination host:port.").Default("127.0.0.1:8356").String()
	)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	kingpin.Version("Ping-pong application")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if *appType == "pong" {

		l, err := net.Listen("tcp", *listenAddress)
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		for {
			// accept connection
			log.Println("Waiting for new connection on address " + *listenAddress + " ...")
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}

			// handle connection
			go HandleConn(conn)
		}

	} else {
		log.Println(*dstAddress)
		//start := time.Now()
		//var m string

		// Get default value for bars
		for i := 0; i < 5; i++ {
			var m string
			start := time.Now()
			m = Ping("tcp", *dstAddress)
			log.Println(i+1, m, " ", time.Now().Format("15:04:05.000000"), " diff: ", time.Since(start).Microseconds())
			stdPong += time.Since(start).Microseconds()
			time.Sleep(200 * time.Millisecond)
		}
		stdPong = int64(float64(stdPong) / 5)
		log.Println(stdPong)

		//ch := make(chan string, 1)
		for i := 0; i < 1000000; i++ {
			var m string
			start := time.Now()
			m = Ping("tcp", *dstAddress)
			since := time.Since(start).Microseconds()
			log.Println(RenderBar(since), m, " diff(microsec): ", since)
			//RenderBar(time.Since(start).Microseconds())
			time.Sleep(200 * time.Millisecond)
		}

	}

}
