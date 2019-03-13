package main

// THIS CODE MAKES NO SENSE -- DO NOT USE IT
// I'm trying to reverse engineer a modbus-based power monitor
// and using it as an excuse to play with Go at the same time.

/* Here's a dump of traffic I captured using a Windows-based application:

2019/03/13 00:01:45  >>> 01 03 00 00 00 0A C5 CD 
2019/03/13 00:01:45  < 01 03 14 00 00 00 00 00 00 00 00 A1 39 63 43 D2 C7 47 42 00 00 00 00 97 4A 
2019/03/13 00:04:38  >>> 01 03 00 00 00 0A C5 CD 
2019/03/13 00:04:38  < 01 03 14 00 00 00 00 00 00 00 00 DF D5 5E 43 DB 94 47 42 00 00 00 00 65 1B 

Anyway, the point of this file isn't to play with a power monitor, it's to
learn Go and get better at programming.

*/

// todo: There are lots of things to do.

import (
	"github.com/goburrow/modbus"
	"time"
	"encoding/binary"
	"fmt"
	"math"
	"github.com/k0kubun/pp"
)

func int16tofloat32(bytearray []byte) float32 {
	buffer := binary.BigEndian.Uint16(bytearray)
	return float32(buffer)
}

func main() {

// Modbus RTU
handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
handler.BaudRate = 9600
handler.DataBits = 8
handler.Parity = "N"
handler.StopBits = 1
handler.SlaveId = 1
handler.Timeout = 5 * time.Second


err := handler.Connect()
defer handler.Close()
if err != nil {
	panic(err)
}

client := modbus.NewClient(handler)

// From the goburrow/modbus/api.go file:
//   results := ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
// Register #5 is what I'm starting with: Voltage.  It should be in the region
// of 200-250 and a 32-bit signed integer.  Modbus uses network byte-order
// (big-endian) and I'm not sure what's going on yet.
//
// There are two 2-word registers (16 bits per register), returning 32 bits as
// I understand it.

results, err := client.ReadHoldingRegisters(5,4)

// The rest of this file is me playing around. This is why it doesn't do
// anything useful.

fmt.Println("Raw output from modbus library")
fmt.Println(results)

fmt.Println("Pretty-Printed:")
pp.Print(results)

fmt.Println("\n--------------------------------")
fmt.Println("Writing into buffer: ")
results_buffer := binary.BigEndian.Uint16(results)
fmt.Println("")
pp.Print(results_buffer)
fmt.Println("")

// golang slice syntax: first value is start address, second is number of bytes
// I think this is where we cut out the addr/opcode from the first two bytes and 
// omit the final CRC, but I'm not sure.
slice_start, slice_bytes := 2,4
float_result := int16tofloat32(results[slice_start:slice_bytes])
fmt.Println(float_result)
fmt.Println("")

bits := binary.BigEndian.Uint32(results)
f := math.Float32frombits(bits)
fmt.Println(f)
}

