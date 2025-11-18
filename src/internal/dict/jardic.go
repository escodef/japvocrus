package dict

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type JardicClient struct {
	baseURL string
	http    *http.Client
}

func NewJardicClient(baseURL string) (*JardicClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("empty Jardic base URL")
	}

	return &JardicClient{
		baseURL: baseURL,
		http:    http.DefaultClient,
	}, nil
}

func (c *JardicClient) GetHTML(word string, page int) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", word)
	q.Set("pg", fmt.Sprintf("%d", page))
	q.Set("sw", "1472")

	u.RawQuery = q.Encode()

	return http.Get(u.String())
}

func (c *JardicClient) GetTranslation(word string) (*Translation, error) {
	resp, err := c.GetHTML(word, 0)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	tab := getElementByID(doc, "tabContent")
	if tab == nil {
		return nil, fmt.Errorf("table with id=tabContent not found")
	}

	firstTR := getFirstChildByTag(tab.LastChild, "tr")
	if firstTR == nil {
		return nil, fmt.Errorf("no translation found")
	}

	td := getChildByTag(firstTR, "td")
	if td == nil {
		return nil, fmt.Errorf("no td found in first translation row")
	}

	reading := getSpanByColor(td, "#7F0000")
	wordJa := getSpanByColor(td, "#00007F")
	translationsText := getSpanByColor(td, "#000000")

	senses := []Sense{}
	for _, line := range strings.Split(translationsText, "<br />") {
		line = strings.TrimSpace(line)
		if line != "" {
			senses = append(senses, Sense{Ru: line})
		}
	}

	return &Translation{
		Word:    wordJa,
		Reading: reading,
		Senses:  senses,
	}, nil
}

func getElementByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := getElementByID(c, id); res != nil {
			return res
		}
	}
	return nil
}

func getFirstChildByTag(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func getChildByTag(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func getSpanByColor(n *html.Node, color string) string {
	var result string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				if a.Key == "style" && strings.Contains(a.Val, "color: "+color) {
					result = getTextContent(n)
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(result)
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(getTextContent(c))
	}
	return sb.String()
}
