package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	targetSites  string
	dork         string
	pages        int
	numThreads   int
	userAgent    string
	searchEngine string
)

func main() {
	displayBanner()
	parseArguments()

	var wg sync.WaitGroup
	urls := make(chan string)

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				fmt.Println("Scanning URL:", url)
			}
		}()
	}

	sites := strings.Split(strings.TrimSpace(targetSites), ",")
	for _, site := range sites {
		for page := 0; page < pages; page++ {
			query := fmt.Sprintf("%s site:%s", dork, site)
			searchURL := constructSearchURL(query, page)
			html, err := fetchHTML(searchURL)
			if err != nil {
				fmt.Printf("Error fetching results for page %d: %v\n", page, err)
				continue
			}
			extractURLs(html, urls)
		}
	}

	close(urls)
	wg.Wait()
}

func displayBanner() {
	blue := color.New(color.FgBlue).SprintFunc()
	banner := blue(`
                                                              .::--=====++++==--::.
                                                        :=+*###***++++++********++++=:
                                                      =******++++++++=++++=====++****++-:
       ..::-==+++++++*++++++++==-:.                 .+********++++******++++**+++=-:::---:
   .-++********+++=+====++++++*****=-.             .+***********#**+****+*****++++++-:
.:-=+******++=====++=++++++++=+*++****+-.       .-++************##***++***********+...
.:---::::::-=+++++++*******+++**++***+==----:---===++*+*******####**+=*********=-. .
        :=++++++**********+++***++*#==---:-+========+++******######+.  .=*#+--:
      ---:.:+#******#****+++**#**+*#***+-:-##+++++===+**#*######*=
           .-:.:++*##***++*#+****+*##***=###***+++++*==+#**##*+
                   .:.=+-.+#+-+***.*#**+:.*******++++*+=*#***:
                                :-    :=+  ***-+****+*+*++:..
                                        .   :-  +*####+*#**+=-.
Dork Hound v1.0.1                             :***#%%##*******+*+=:.
Coded by: kaizer baynosa | spectre            -*#*##*%###**+**+=:..::--:.
                                             +**#*#*######***++***+-.
                                            *#-+**#**#+**##*+**+=-:-=+=:
                                           +*. ***+****+**##*+-+*+-    .---:
                                           .. .**+ **#*++*+*#*+.:**+:
                                              =**. **#.*+**+***+- -**=.
                                              **:  =**  ++****-**=  .=*=.
                                             .#-   -**  .+*#**= =*+.    -+:
                                             :+    :**   .+**+*- :**:     :=:
                                                    *+    -*=++*-  =*+       .
                                                    *+     +* =+*.   -*:
                                                    ==      *= -++     -=.
                                                    ::      .*- .=.      ::
                                                              +-
`)
	fmt.Println(banner)
}

func parseArguments() {
	flag.StringVar(&targetSites, "u", "", "Target sites (comma-separated)")
	flag.StringVar(&dork, "d", "", "Google dorking query")
	flag.IntVar(&pages, "p", 0, "Number of pages to search")
	flag.IntVar(&numThreads, "t", 0, "Number of threads")
	flag.StringVar(&userAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3", "User-Agent header")
	flag.StringVar(&searchEngine, "e", "Google", "Search engine choice (Google, Bing, DuckDuckGo)")
	flag.Parse()
}

func constructSearchURL(query string, page int) string {
	var baseSearchURL string
	switch searchEngine {
	case "Google":
		baseSearchURL = "https://www.google.com/search"
	case "Bing":
		baseSearchURL = "https://www.bing.com/search"
	case "DuckDuckGo":
		baseSearchURL = "https://duckduckgo.com/html"
	default:
		baseSearchURL = "https://www.google.com/search"
	}
	return fmt.Sprintf("%s?q=%s&start=%d", baseSearchURL, query, page*10)
}

func fetchHTML(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	var htmlBuilder strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 || err != nil {
			break
		}
		htmlBuilder.Write(buf[:n])
	}

	return htmlBuilder.String(), nil
}

func extractURLs(html string, urls chan<- string) {
	re := regexp.MustCompile(`<a\s+href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		url := match[1]
		if strings.HasPrefix(url, "http") {
			urls <- url
		}
	}
}
//hello 
