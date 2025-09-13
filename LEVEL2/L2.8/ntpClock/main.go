// Program ntpClock fetches and prints the current time from an NTP server.
package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	srvTime, err := ntp.Time("0.pool.ntp.org")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to fetch time:", err)
		os.Exit(1)
	}

	fmt.Println("UTC-timezone time:", srvTime.UTC())
	fmt.Println("Local timezone time:", srvTime.Local())
}
