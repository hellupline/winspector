package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Record struct {
	binKey     uuid.UUID
	recordKey  uuid.UUID
	createdAt  time.Time
	recordData RecordData
}

func NewRecord(binKey, recordKey uuid.UUID, createdAt time.Time, recordData RecordData) Record {
	return Record{
		binKey:     binKey,
		recordKey:  recordKey,
		createdAt:  createdAt,
		recordData: recordData,
	}
}

type RecordData struct {
	method           string
	uRL              string
	proto            string
	host             string
	transferEncoding []string
	contentLength    int64
	headers          PairList
	query            PairList
	postFormData     PairList
	body             []byte
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	if p[i].key == p[j].key {
		return p[i].value < p[j].value
	}
	return p[i].key < p[j].key
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Pair struct {
	key   string
	value string
}

func NewRecordData(r *http.Request) RecordData {
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
			break
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
	return RecordData{
		method:           r.Method,
		uRL:              r.URL.String(),
		proto:            r.Proto,
		host:             r.Host,
		transferEncoding: r.TransferEncoding,
		contentLength:    r.ContentLength,
		headers:          headers,
		query:            query,
		postFormData:     postFormData,
		body:             body,
	}
}

var (
	headersIgnoreList = map[string]bool{
		"x-forwarded-for":    true,
		"x-forwarded-host":   true,
		"x-forwarded-port":   true,
		"x-forwarded-proto":  true,
		"x-forwarded-server": true,
		"x-real-ip":          true,
	}
)
