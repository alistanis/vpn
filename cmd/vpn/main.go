package main

import (
	"os"

	"github.com/alistanis/vpn"
)

func main() {
	// this returns an error but cobra prints it anyway
	if err := vpn.ExecuteRoot(); err != nil {
		os.Exit(1)
	}
}
