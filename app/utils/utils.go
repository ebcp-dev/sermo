package app

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// Error message response.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// JSON http response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Signalling server helpers.

// Starts HTTP server that consumes SDPs.
func HTTPSDPServer() chan string {
	port := flag.Int("port", viper.GetInt("SDP_PORT"), "http server port")
	flag.Parse()

	sdpChan := make(chan string)
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		log.Println(w, "done")
		sdpChan <- string(body)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
		if err != nil {
			panic(err)
		}
	}()

	return sdpChan
}

// Blocks until input is received from stdin.
func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)
	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	log.Println("")

	return in
}

// Encodes input in base64.
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	// Optional zip before encoding.
	if compress {
		b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decodes input from base64.
func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	// Optional zip before decoding.
	if compress {
		b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

// Compress input string.
func zip(in []byte) []byte {
	var b bytes.Buffer

	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}

	err = gz.Flush()
	if err != nil {
		panic(err)
	}

	err = gz.Close()
	if err != nil {
		panic(err)
	}

	return b.Bytes()
}

// Decompress input string.
func unzip(in []byte) []byte {
	var b bytes.Buffer

	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}

	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}

	res, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return res
}
