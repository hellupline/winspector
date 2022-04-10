package models

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sort"
)

type RequestData struct {
	Method           string
	URL              string
	Proto            string
	Host             string
	RemoteAddr       string
	RequestURI       string
	TransferEncoding []string
	ContentLength    int64
	Headers          PairList
	Query            PairList
	PostFormData     PairList
	Body             []byte
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	if p[i].Key == p[j].Key {
		return p[i].Value < p[j].Value
	}
	return p[i].Key < p[j].Key
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Pair struct {
	Key   string
	Value string
}

func NewRequestData(r *http.Request) RequestData {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println(err)
	}
	headers := PairList{}
	for k, vArr := range r.Header {
		if _, ok := headersIgnoreList[k]; ok {
			continue
		}
		for _, v := range vArr {
			headers = append(headers, Pair{k, v})
		}
	}
	sort.Sort(headers)
	query := PairList{}
	for k, vArr := range r.URL.Query() {
		for _, v := range vArr {
			query = append(query, Pair{k, v})
		}
	}
	sort.Sort(query)
	postFormData := PairList{}
	for k, vArr := range r.PostForm {
		for _, v := range vArr {
			postFormData = append(postFormData, Pair{k, v})
		}
	}
	sort.Sort(postFormData)
	return RequestData{
		Method:           r.Method,
		URL:              r.URL.String(),
		Proto:            r.Proto,
		Host:             r.Host,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TransferEncoding: r.TransferEncoding,
		ContentLength:    r.ContentLength,
		Headers:          headers,
		Query:            query,
		PostFormData:     postFormData,
		Body:             body,
	}
}

var headersIgnoreList = map[string]bool{
	"X-Forwarded-For":    true,
	"X-Forwarded-Host":   true,
	"X-Forwarded-Port":   true,
	"X-Forwarded-Proto":  true,
	"X-Forwarded-Server": true,
	"X-Real-Ip":          true,
}
