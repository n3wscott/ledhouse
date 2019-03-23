package main

import "fmt"
import "time"
import "log"
import "github.com/jacobsa/go-serial/serial"
import "flag"
import "os"
import "io"

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

var options struct {
	SerialPort string
	SerialBaud uint
}

func init() {
	flag.StringVar(&options.SerialPort, "serial", "", "The serial port to use to connect to the ledhouse.")
	flag.UintVar(&options.SerialBaud, "baud", 19200, "The baud rate the serial port will use.")
	flag.Parse()
}

func main() {

	fmt.Println("port:", options.SerialPort)
	//	portName := flag.String("port", "", "serial port to test (/dev/ttyUSB0, etc)")

	if options.SerialPort == "" {
		fmt.Println("Must specify port")
		usage()
	}

	// Set up options.
	options := serial.OpenOptions{
		PortName:        options.SerialPort,
		BaudRate:        options.SerialBaud,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

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

	// Make sure to close it later.
	defer port.Close()

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

	time.Sleep(time.Millisecond * 150)

	{
		_, err := port.Write([]byte("CLEAR\n"))
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}
		fmt.Println("CLEAR")
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
