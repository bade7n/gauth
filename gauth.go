package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/bade7n/gauth/gauth"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	cfgPath := os.Getenv("GAUTH_CONFIG")
	if cfgPath == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		cfgPath = path.Join(user.HomeDir, ".config/gauth.csv")
	}

	cfgContent, err := gauth.LoadConfigFile(cfgPath, getPassword)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	urls, err := gauth.ParseConfig(cfgContent)
	if err != nil {
		log.Fatalf("Decoding configuration file: %v", err)
	}

	_, progress := gauth.IndexNow() // TODO: do this per-code

	if len(os.Args) > 1 {
		account := os.Args[1]
		for _, url := range urls {
			if account == url.Account {
				_, curr, _, _ := gauth.Codes(url)
				fmt.Fprintf(os.Stdout, "%s\n", curr)
			}
		}

	} else {
		tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)
		fmt.Fprintln(tw, "\tprev\tcurr\tnext")
		for _, url := range urls {
			prev, curr, next, err := gauth.Codes(url)
			if err != nil {
				log.Fatalf("Generating codes for %q: %v", url.Account, err)
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", url.Account, prev, curr, next)
		}
		tw.Flush()
		fmt.Printf("[%-29s]\n", strings.Repeat("=", progress))

	}

}

func getPassword() ([]byte, error) {
	fmt.Printf("Encryption password: ")
	defer fmt.Println()
	return terminal.ReadPassword(int(syscall.Stdin))
}
