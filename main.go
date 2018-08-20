package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var c *cache
var port *string
var redirecturl *string
var configFile *string
var cacheHTTPHeader *string
var active bool

type prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

type transport struct {
	http.RoundTripper
}

type cache struct {
	mimeTypes []string
	pages     []string
	responses map[string]cachedResponse
}

type cachedResponse struct {
	body   []byte
	header map[string][]string
}

func main() {
	const (
		defaultPort            = "80"
		defaultTarget          = "http://localhost:8080"
		defaultConfigFile      = "config/cache.config"
		defaultcacheHTTPHeader = "X-GoProxy"
	)

	// flags and env variables
	envP, ex := os.LookupEnv("PROXY_PORT")
	if ex {
		port = flag.String("port", envP, "default server port, ':"+envP+"'")
	} else {
		port = flag.String("port", defaultPort, "default server port, ':"+defaultPort+"'")
	}
	envT, ex := os.LookupEnv("PROXY_TARGET")
	if ex {
		redirecturl = flag.String("url", envT, "default redirect url, '"+envT+"'")
	} else {
		redirecturl = flag.String("url", defaultTarget, "default redirect url, '"+defaultTarget+"'")
	}
	envC, ex := os.LookupEnv("CONFIG_FILE")
	if ex {
		configFile = flag.String("configFile", envC, "default config file, '"+envC+"'")
	} else {
		configFile = flag.String("configFile", defaultConfigFile, "default config file, '"+defaultConfigFile+"'")
	}
	envH, ex := os.LookupEnv("PROXY_CACHE_HTTP_HEADER")
	if ex {
		cacheHTTPHeader = flag.String("cacheHttpHeader", envH, "default config file, '"+envH+"'")
	} else {
		cacheHTTPHeader = flag.String("cacheHttpHeader", defaultcacheHTTPHeader, "default config file, '"+defaultcacheHTTPHeader+"'")
	}

	flag.Parse()

	fmt.Println("server will run on :", *port)
	fmt.Println("redirecting to :", *redirecturl)

	//cache
	c = newCache(configFile)

	//cache active
	active = true

	// proxy
	proxy := newProxy(*redirecturl)

	http.HandleFunc("/healthcheck", healthcheck)

	http.HandleFunc("/flushcache", flushcache)

	http.HandleFunc("/cacheactive", cacheactive)

	// server redirection
	http.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func newCache(configFile *string) *cache {
	content, err := ioutil.ReadFile(*configFile)
	if err != nil {
		//Do something
	}
	cleanLines := strings.Replace(string(content), "\r", "", -1)
	var mimeTypes []string
	var pages []string
	for _, str := range strings.Split(cleanLines, "\n") {
		if str != "" {
			if strings.Contains(str, "{page}") {
				pages = append(pages, strings.Replace(str, "{page}", "", -1))
			} else if strings.Contains(str, "{mimetype}") {
				mimeTypes = append(mimeTypes, strings.Replace(str, "{mimetype}", "", -1))
			}
		}
	}
	return &cache{responses: map[string]cachedResponse{}, mimeTypes: mimeTypes, pages: pages}
}

func newProxy(target string) *prox {
	url, _ := url.Parse(target)

	return &prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Port:" + *port + ", Redirecting to :" + *redirecturl))
}

func flushcache(w http.ResponseWriter, r *http.Request) {
	c = newCache(configFile)
	w.Write([]byte("Cache Flushed"))
}

func cacheactive(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		active = !active
		if active {
			c = newCache(configFile)
		}
	}

	var activeText = "Deactivate"

	if !active {
		activeText = "Activate"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<form action=\"/cacheactive\" method=\"post\"><button type=\"submit\">" + activeText + "</button></form>"))
}

func (p *prox) handle(w http.ResponseWriter, r *http.Request) {

	p.proxy.Transport = &transport{http.DefaultTransport}

	p.proxy.ServeHTTP(w, r)
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if active && req.Method == "GET" {

		f, ok := c.responses[req.RequestURI]

		if ok {

			body := ioutil.NopCloser(bytes.NewReader(f.body))

			resp = &http.Response{Body: body, Header: f.header, StatusCode: 200}

			resp.Header.Set(*cacheHTTPHeader, "FromCache")

			resp.ContentLength = int64(len(f.body))
			return resp, nil
		}
	}

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if active && req.Method == "GET" && resp.StatusCode == 200 {

		var isC bool

		cType := resp.Header.Get("content-type")

		for _, a := range c.mimeTypes {
			if strings.Contains(cType, a) {
				isC = true
				break
			}
		}

		if !isC {
			for _, str := range c.pages {
				if req.RequestURI == str {
					isC = true
					break
				}
			}
		}

		if isC {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = resp.Body.Close()
			if err != nil {
				return nil, err
			}

			body := ioutil.NopCloser(bytes.NewReader(b))
			resp.Body = body

			c.responses[req.RequestURI] = cachedResponse{body: b, header: resp.Header}
		}
	}

	return resp, nil
}
