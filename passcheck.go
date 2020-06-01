package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func haveIBeenPwnedPasswordsLookup(password string) int {
	hash := sha1.New()
	hash.Write([]byte(password))
	sha1_hash := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	hash_head := sha1_hash[:5]
	hash_tail := sha1_hash[5:]

	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + hash_head)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(body), "\r\n")
	for _, line := range lines {
		line_parts := strings.Split(line, ":")
		tail := line_parts[0]
		if tail == hash_tail {
			passwordCount, err := strconv.Atoi(line_parts[1])
			if err != nil {
				fmt.Println(err)
				return -1
			}
			return passwordCount
		}
	}
	return -1
}

func customUsage() {
	fmt.Fprintf(os.Stderr,
		`Usage: %s <flags>

Checks passwords against haveibeenpwned.com/passwords

Flags:
`,
		os.Args[0])

	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage

	var workers int
	flag.IntVar(&workers, "w", 20, "the ammount of workers")

	var showPasswordCount bool
	flag.BoolVar(&showPasswordCount, "c", false, "show password count")

	flag.Parse()
	passwords := make(chan string)
	output := make(chan string)

	// have I been pwned worker
	var hibpWG sync.WaitGroup
	for i := 0; i < workers; i++ {
		hibpWG.Add(1)

		go func() {
			for password := range passwords {
				if passwordCount := haveIBeenPwnedPasswordsLookup(password); passwordCount > 0 {
					if showPasswordCount {
						output <- fmt.Sprintf("%s:%d", password, passwordCount)
					} else {
						output <- password
					}
					continue
				}
			}
			hibpWG.Done()
		}()

	}

	// output worker
	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		for o := range output {
			fmt.Println(o)
		}
		outputWG.Done()
	}()

	// close the output channel when have I been pwned workers are done
	go func() {
		hibpWG.Wait()
		close(output)
	}()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		password := sc.Text()
		passwords <- password
	}

	close(passwords)

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}
	outputWG.Wait()

}
