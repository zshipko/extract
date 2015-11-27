package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var urls []string = make([]string, 0)
var tags []string = make([]string, 0)
var classes []string = make([]string, 0)
var ids []string = make([]string, 0)
var search string = ""
var printed []string = make([]string, 0)

func printl(s string) {
	s = strings.TrimSpace(s)
	if len(s) > 0 && !in(printed, s) {
		printed = append(printed, s)
		fmt.Println(s)
	}
}

func in(hay []string, needle string) bool {
	for _, x := range hay {
		if x == needle {
			return true
		}
	}
	return false
}

func fixUrl(home, u string) string {
	h, err := url.Parse(home)
	if err != nil {
		return ""
	}

	if strings.HasPrefix(u, "http") {
		return u
	} else if strings.HasPrefix(u, "//") {
		return "http:" + u
	} else if strings.HasPrefix(u, "/") {
		if len(h.Scheme) == 0 {
			return "http://" + h.Host + u
		}
		return h.Scheme + "://" + h.Host + u
	} else if strings.HasPrefix(u, "www.") {
		return "http://" + u
	} else if strings.HasPrefix(u, "mailto:") {
		return strings.Split(u, ":")[1]
	}
	return path.Join(home, u)
}

func parseHtml(u string, r io.Reader) {
	d := html.NewTokenizer(r)
	for {
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := d.Token()
		switch tokenType {
		case html.StartTagToken: // <tag>
			if !in(tags, token.Data) {
				continue
			}
			switch token.Data {
			case "a":
				for _, a := range token.Attr {
					if a.Key == "href" &&
						strings.Contains(a.Val, search) &&
						(!strings.HasPrefix(a.Val, "mailto:") || in(tags, "email")) {
						if !strings.HasPrefix(a.Val, "mailto:") && in(tags, "email") {
							continue
						}
						printl(fixUrl(u, a.Val) + " ")
					}
				}
			case "img":
				for _, a := range token.Attr {
					if a.Key == "src" && strings.Contains(a.Val, search) {
						printl(fixUrl(u, a.Val))
					}
				}
			}

		}
	}
}

func getPage(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	parseHtml(u, resp.Body)
	return nil
}

func main() {
	argc := len(os.Args)
	for i := 1; i < argc; i++ {
		arg := os.Args[i]
		if arg[0] == '-' {
			tags = append(tags, arg[1:])
		} else if arg[0] == '?' {
			search = arg[1:]
		} else {
			if !strings.HasPrefix(arg, "http") {
				arg = "http://" + arg
			}
			urls = append(urls, arg)
		}
	}

	if len(urls) == 0 {
		parseHtml("/", os.Stdin)
	}

	if len(tags) == 0 {
		tags = append(tags, "img")
		tags = append(tags, "a")
	}

	for _, u := range urls {
		getPage(u)
	}
}
