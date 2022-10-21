package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/iangcarroll/cookiemonster/pkg/monster"
)

var (
	concurrencyFlag = flag.Int("concurrency", runtime.NumCPU(), "Optional. How many attempts should run concurrently; the default is runtime.NumCPU().")
	cookiesFlag     = flag.String("cookies", "", "Required. Path to load cookie values from.")
	wordlistFlag    = flag.String("wordlist", "", "Required. Path to load a base64-encoded wordlist from.")
)

func MonsterRun(cookie string, wl *monster.Wordlist) (success bool, err error) {
	c := monster.NewCookie(cookie)

	if !c.Decode() {
		return false, errors.New("couldNotDecode")
	}

	if _, success := c.Unsign(wl, uint64(*concurrencyFlag)); !success {
		return false, errors.New("couldNotUnsign")
	}

	return true, nil
}

func main() {
	flag.Parse()

	if *cookiesFlag == "" || *wordlistFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// fetch wordlist/secrets
	wl := monster.NewWordlist()

	if err := wl.Load(*wordlistFlag); err != nil {
		fmt.Printf("Sorry, I could not load your wordlist. Please ensure every line contains valid base64. Error: %v", err)
		os.Exit(1)
	}

	// loop over cookie values
	cookies, err := os.Open(*cookiesFlag)

	if err != nil {
		fmt.Printf("Sorry, I could not load your cookies. Error: %v", err)
		os.Exit(1)
	}
	defer cookies.Close()

	scanner := bufio.NewScanner(cookies)
	for scanner.Scan() {
		success, err := MonsterRun(scanner.Text(), wl)

		if !success {
			fmt.Println(err, scanner.Text())
		} else {
			fmt.Println("success", scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
