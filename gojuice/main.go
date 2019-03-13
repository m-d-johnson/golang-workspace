package main

import (
	"github.com/goburrow/modbus"
	"time"
	"encoding/binary"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/google/gopacket"
)

func int16tofloat32(bytearray []byte) float32 {
	buffer := binary.BigEndian.Uint16(bytearray)
	return float32(buffer)
}

func main() {

// Modbus RTU/ASCII
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
//results := ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)

results, err := client.ReadHoldingRegisters(1,20)
fmt.Println("Raw output from modbus library")
fmt.Println(results)
fmt.Println("Pretty-Printed:")
pp.Print(results)

fmt.Println("Writing into buffer as BigEndian")
results_buffer := binary.BigEndian.Uint16(results)
pp.Print(results_buffer)

// golang slice syntax: first value is start address, second is number of bytes
// I think this is where we cut out the addr/opcode from the first two bytes and 
// omit the final CRC, but I'm not sure.
slice_start, slice_bytes := 0,4
float_result := int16tofloat32(results[slice_start:slice_bytes])
fmt.Printf("%f\n\n", float_result)

}

