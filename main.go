package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	base64TranslationTable = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	MASK_6_BIT             = uint32(0x3F) // 0x3F = 00000000 00000000 00000000 00111111
)

func debug(s string) {
	if os.Getenv("DEBUG") == "1" {
		fmt.Println(s)
	}
}

func indexOfChar(s string) int {
	for i := 0; i < len(base64TranslationTable); i++ {
		if s == string(base64TranslationTable[i]) {
			return i
		}
	}

	return 0
}

func outputEncodedData(writer *bufio.Writer, data uint32) {
	writer.WriteString(string(base64TranslationTable[data]))
}

func outputDecodedData(writer *bufio.Writer, data byte) {
	writer.WriteByte(data)
}

func encodeToBase64(input io.Reader, output io.Writer) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	defer writer.Flush()

	byteChunkSize := 3 // 24 bits
	readbuf := make([]byte, byteChunkSize)
	// chunks := []uint32{}
	combinedbyts := uint32(0) // 00000000 00000000 00000000 00000000
	bit_count := 0

	padding := 0
	for {

		bytesReadCount, err := reader.Read(readbuf)

		if bytesReadCount > 0 {
			debug(fmt.Sprintf("%08b %d", readbuf, bytesReadCount))
			for i := range bytesReadCount {

				combinedbyts = (combinedbyts << 8) | uint32(readbuf[i])

				debug(fmt.Sprintf("combinedbyts = %032b", combinedbyts))

				bit_count += 8
				debug(fmt.Sprintf("bit_count=%d", bit_count))
				for bit_count >= 6 {
					bit_count -= 6
					// get first 6 bits from left , disgard rest bits
					chunk := (combinedbyts >> bit_count) & MASK_6_BIT

					debug(fmt.Sprintf("%032b \n%032b \n%032b", combinedbyts>>bit_count, MASK_6_BIT, chunk))

					// chunks = append(chunks, chunk)
					outputEncodedData(writer, chunk)

					debug("--------------------")
				}

			}

		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panic(err)
		}

	}

	if bit_count > 0 {
		debug(fmt.Sprintf("bit_count=%d", bit_count))
		trailing_bits := 6 - bit_count
		debug(fmt.Sprintf("trailing bit_count=%d", trailing_bits))
		chunk := (combinedbyts << trailing_bits) & MASK_6_BIT

		debug(fmt.Sprintf("%032b \n%032b \n%032b \n%032b", combinedbyts, combinedbyts<<trailing_bits, MASK_6_BIT, chunk))

		// chunks = append(chunks, chunk)
		outputEncodedData(writer, chunk)

		debug("--------------------")

		padding += trailing_bits / 2

	}

	for i := 0; i < padding; i++ {
		// chunks = append(chunks, 64)
		outputEncodedData(writer, 64)
	}
	// outputEncodedData(chunks)
}

func decodeFromBase64(input io.Reader, ouput io.Writer) {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(ouput)

	defer writer.Flush()

	byteChunkSize := 4096
	readbuf := make([]byte, byteChunkSize)
	// chunks := []byte{}
	combinedbyts := uint32(0) // 00000000 00000000 00000000 00000000
	bit_count := 0

	for {
		bytesReadCount, err := reader.Read(readbuf)

		if bytesReadCount > 0 {
			debug(fmt.Sprintf("%s", readbuf))
			for i := range bytesReadCount {
				// fmt.Printf("newline = %c %b %d %08x\n", readbuf[i], readbuf[i], readbuf[i], readbuf[i])
				// skip new line
				if readbuf[i] == '\n' || readbuf[i] == '\r' {
					continue
				}
				c := indexOfChar(string(readbuf[i]))
				debug(fmt.Sprintf("%c %08b %d -> %d %08b", readbuf[i], readbuf[i], readbuf[i], c, c))
				combinedbyts = (combinedbyts << 6) | (uint32(c) & MASK_6_BIT)
				debug(fmt.Sprintf("%032b", combinedbyts))
				bit_count += 6

				for bit_count >= 24 {
					debug(fmt.Sprintf("bit_count = %d", bit_count))
					bit_count -= 8
					chunk := byte(combinedbyts>>bit_count) & 0xFF
					debug(fmt.Sprintf("chunk --------> %08b", chunk))

					// chunks = append(chunks, chunk)
					outputDecodedData(writer, chunk)
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panic(err)
		}

	}

	for bit_count > 0 {
		debug(fmt.Sprintf("last bit_count = %d", bit_count))
		bit_count -= 8
		// 00000000010011010110000101101110
		chunk := byte(combinedbyts>>bit_count) & 0b11111111
		debug(fmt.Sprintf("chunk --------> %08b %d", chunk, chunk))
		if chunk == 0 {
			continue
		}
		// chunks = append(chunks, chunk)
		outputDecodedData(writer, chunk)
	}
	// debug(fmt.Sprintf("%d", chunks))
	// outputDecodedData(chunks)
}

func base64(input io.Reader, outout io.Writer, decode bool) {
	if decode {
		decodeFromBase64(input, outout)
	} else {
		encodeToBase64(input, outout)
		println()
	}
}

func main() {
	args := os.Args
	decode := false
	decodeFlag := "-d"

	if len(args) > 1 && args[1] == decodeFlag {
		decode = true
	}

	input := os.Stdin
	output := os.Stdout

	if len(args) > 1 {
		inputfile := 1
		if args[inputfile] == decodeFlag {
			inputfile += 1
		}
		if inputfile == (len(args) - 1) {

			fd, err := os.Open(args[inputfile])
			defer fd.Close()
			if err != nil {
				log.Panic(err)
			}
			input = fd
		}
	}

	base64(input, output, decode)
}
