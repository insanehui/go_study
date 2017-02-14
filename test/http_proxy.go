// package main

// import (
// 	"net/url"
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"os"
// )

// func main() {
// 	targetUrl, err := url.Parse("http://www.baidu.com")

// 	if err != nil {
// 		panic("bad url")
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

// 	http.Handle("/", proxy)

// 	log.Println("Start serving on port 1234")

// 	http.ListenAndServe(":1234", nil)

// 	os.Exit(0)
// }

// package main

// import (
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// )

// func main() {
// 	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
// 		Scheme: "http",
// 		Host:   "www.baidu.com",
// 	})
// 	http.ListenAndServe(":9090", proxy)
// }

package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var targetURL *url.URL

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func handler(w http.ResponseWriter, r *http.Request) {
	o := new(http.Request)

	*o = *r

	o.Host = targetURL.Host
	o.URL.Scheme = targetURL.Scheme
	o.URL.Host = targetURL.Host
	o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)

	if q := o.URL.RawQuery; q != "" {
		o.URL.RawPath = o.URL.Path + "?" + q
	} else {
		o.URL.RawPath = o.URL.Path
	}

	o.URL.RawQuery = targetURL.RawQuery

	o.Proto = "HTTP/1.1"
	o.ProtoMajor = 1
	o.ProtoMinor = 1
	o.Close = false

	transport := http.DefaultTransport

	res, err := transport.RoundTrip(o)

	if err != nil {
		log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hdr := w.Header()

	for k, vv := range res.Header {
		for _, v := range vv {
			hdr.Add(k, v)
		}
	}

	for _, c := range res.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}

	w.WriteHeader(res.StatusCode)

	if res.Body != nil {
		io.Copy(w, res.Body)
	}
}

func main() {
	url, err := url.Parse("http://www.baidu.com")

	if err != nil {
		log.Println("Bad target URL")
	}

	targetURL = url

	http.HandleFunc("/", handler)

	log.Println("Start serving on port 1234")

	http.ListenAndServe(":1234", nil)
}
