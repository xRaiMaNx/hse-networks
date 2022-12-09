package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

const (
	headersSize = 28 // in bytes
)

func tryPing(dest string, count int, v bool) (int, error) {
	cmd := exec.Command("ping", dest, "-c", "1", "-s", strconv.Itoa(count), "-M", "do")
	out, err := cmd.Output()
	if v {
		log.Println(string(out))
	}
	if werr, ok := err.(*exec.ExitError); ok {
		code := werr.ExitCode()
		return code, err
	}
	return 0, err
}

func main() {
	host := flag.String("host", "", "destination for ping")
	v := flag.Bool("verbose", false, "verbose log")
	flag.Parse()

	right := 1

	// check correctness of destination or other errors
	code, err := tryPing(*host, right, *v)
	if code != 0 {
		log.Fatalf("failed with code: %v\terror: %v", code, err)
	}

	// find the upper threshold O(logN) where N - MTU
	// default is 1500 bytes, but it can be bigger. For instance, fast Ethernet links
	left := right // for optimizing future bin search
	for code == 0 {
		left = right
		right *= 2
		code, err = tryPing(*host, right, *v)
		if code == 2 {
			log.Fatalf("failed with code: %v\terror: %v", code, err)
		}
	}
	if *v {
		log.Printf("find the upper threshold: %v\n", right)
	}

	for right > left+1 {
		mid := (right + left) / 2
		code, err = tryPing(*host, mid, *v)
		if *v {
			log.Println(left, right)
		}
		if code == 2 {
			log.Fatalf("failed with code: %v\terror: %v", code, err)
		}
		if code == 1 {
			right = mid
		} else {
			left = mid
		}
	}

	fmt.Printf("MTU is %v\n", left+headersSize)
}
