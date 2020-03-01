package main

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"time"
 "log"
"github.com/pkg/term"
 "flag"
 "os"
"io"
)

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

var options struct {
	SerialPort string
	SerialBaud int
	Loops int
}

func init() {
	flag.StringVar(&options.SerialPort, "serial", "", "The serial port to use to connect to the ledhouse.")
	flag.IntVar(&options.SerialBaud, "baud", 19200, "The baud rate the serial port will use.")
	flag.IntVar(&options.Loops, "loops", 1, "The number of loops to perform. 0 means forever.")
	flag.Parse()
}

func main() {

	fmt.Println("port:", options.SerialPort)
	fmt.Println("loops:", options.Loops)

	if options.SerialPort == "" {
		fmt.Println("Must specify port")
		usage()
	}

	t, err := term.Open(options.SerialPort)
	if err != nil {
		log.Fatalf("term.Open: %v", err)
	}

	if err := t.SetSpeed(options.SerialBaud); err != nil {
		log.Fatalf("t.SetSpeed: %v", err)
	}



	//// Set up options.
	//opts := serial.OpenOptions{
	//	PortName:        options.SerialPort,
	//	BaudRate:        options.SerialBaud,
	//	DataBits:        8,
	//	StopBits:        1,
	//	MinimumReadSize: 4,
	//}
	//
	//// Open the port.
	//port, err := serial.Open(opts)
	//if err != nil {
	//	log.Fatalf("serial.Open: %v", err)
	//}
	//// Make sure to close it later.
	//defer port.Close()

	{
		_, err := t.Write([]byte("\n\n\n"))
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}
	}

	if options.Loops == 0 {
		for {
			loop(t)
		}
	} else {
		for i := 0; i < options.Loops; i++ {
			loop(t)
		}
	}

	{
		time.Sleep(time.Millisecond * 150)
		_, err := t.Write([]byte("CLEAR\n"))
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}
		fmt.Println("CLEAR")
	}
}

func loop(port *term.Term) {
	{
		buf := make([]byte, 5)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			fmt.Println("Rx: ", buf)
		}
	}

	for offset := 0; offset < 100; offset++ {
		for i := 0; i < 10; i++ {
			color := Wheel((15 * (i + offset)) % 255)
			cmd := fmt.Sprintf("LED%02x%s\n", i, color)
			_, err := port.Write([]byte(cmd))
			if err != nil {
				log.Fatalf("port.Write: %v", err)
			}
			fmt.Print(cmd)
		}
		time.Sleep(time.Millisecond * 150)
	}
}

// Adapted from the neopixel example code.
func Wheel(pos int) string {
	pos = 255 - pos
	var r, g, b int
	switch {
	case pos < 85:
		r = 255 - pos*3
		g = 0
		b = pos * 3
	case pos < 170:
		pos -= 85
		r = 0
		g = pos * 3
		b = 255 - pos*3
	default:
		pos -= 170
		r = pos * 3
		g = 255 - pos*3
		b = 0
	}
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}
