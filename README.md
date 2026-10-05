# Custom Base64 Encoder & Decoder (Go)

A lightweight, from-scratch implementation of Base64 encoding and decoding in Go using low-level bitwise operations as specified in [RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648#section-4).

> **Note**: This project was built for **educational and learning purposes** to explore and understand how Base64 bit manipulation, buffering, and padding work under the hood without relying on Go's built-in `encoding/base64` package. It was implemented while referencing Wikipedia and YouTube video explanations, so the code is unoptimized and messy. It is not intended for production use.

---

## Features

- **Standard Base64 [RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648#section-4)**: Uses the standard 64-character alphabet (`A-Z`, `a-z`, `0-9`, `+`, `/`) and `=` padding.
- **Bitwise Stream Processing**: Encodes/decodes by bit-shifting byte streams into 6-bit and 8-bit chunks.
- **Flexible Input**: Supports reading from **Standard Input (stdin)** or from a **file path**.
- **Interactive Debug Mode**: Inspect binary bit representations and buffer states step-by-step using an environment variable.

---

## How It Works (Learning Concepts)

1. **Encoding**:
   - Reads input bytes (8 bits each).
   - Combines bytes into a buffer and extracts 6-bit chunks (`MASK_6_BIT = 0x3F`).
   - Maps each 6-bit integer (0–63) to its corresponding character in the Base64 index table.
   - Adds `=` padding characters when the total number of bits is not a multiple of 24 (3 bytes).

2. **Decoding**:
   - Reads Base64 characters and converts each character back to its 6-bit value.
   - Combines 6-bit values and extracts 8-bit bytes (`0xFF`).
   - Writes the reconstructed original bytes to standard output.

---

## Prerequisites

- [Go](https://go.dev/dl/) 1.20+ installed.

---

## Usage

You can run the program directly with `go run main.go` or build a binary executable.

### 1. Build the Binary (Optional)

```bash
go build -o base64 main.go
```

### 2. Encoding to Base64

**From a file:**
```bash
go run main.go test_files/hello.txt
# or with compiled binary:
./base64 test_files/hello.txt
```

**From standard input (stdin):**
```bash
echo -n "Hello, World!" | go run main.go
```

---

### 3. Decoding from Base64

Use the `-d` flag to decode Base64 data back to plain text/binary.

**From a file:**
```bash
go run main.go -d test_files/hello.txt.b64
# or with compiled binary:
./base64 -d test_files/hello.txt.b64
```

**From standard input (stdin):**
```bash
echo "SGVsbG8sIFdvcmxkIQ==" | go run main.go -d
```

---

### 4. Debug Mode

To see binary representations, intermediate bit counts, and buffer shifts in real-time, set `DEBUG=1`:

```bash
DEBUG=1 go run main.go test_files/hello.txt
```

---

## Running Tests

Run the test suite using Go's built-in testing tool:

```bash
go test -v ./...
```

---

## Project Structure

```text
.
├── .gitignore         # Git ignore rules
├── LICENSE            # MIT License
├── README.md          # Project documentation
├── go.mod             # Go module definition
├── main.go            # Custom Base64 encoding & decoding logic
└── main_test.go       # Unit tests for Base64 encode and decode
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Disclaimer

This repository is intended solely for **learning and demonstration purposes**. For production environments, use Go's standard library package [`encoding/base64`](https://pkg.go.dev/encoding/base64) which includes optimizations and SIMD accelerations.
