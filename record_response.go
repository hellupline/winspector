package main

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type RecordResponseList []RecordResponse

func (o RecordResponseList) Len() int {
	return len(o)
}

func (o RecordResponseList) Less(i, j int) bool {
	return o[i].CreatedAt > o[j].CreatedAt
}

func (o RecordResponseList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func NewRecordResponseList(recordStore map[uuid.UUID]Record) RecordResponseList {
	recordResponseList := make(RecordResponseList, 0, len(recordStore))
	for _, record := range recordStore {
		recordResponse := NewRecordResponse(record)
		recordResponseList = append(recordResponseList, recordResponse)
	}
	sort.Sort(recordResponseList)
	return recordResponseList
}

type PairResponseList []PairResponse

type PairResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RecordResponse struct {
	BinKey           string           `json:"bin_key"`
	RecordKey        string           `json:"record_key"`
	CreatedAt        string           `json:"created_at"`
	Method           string           `json:"method"`
	URL              string           `json:"url"`
	Proto            string           `json:"proto"`
	Host             string           `json:"host"`
	RemoteAddr       string           `json:"remote_addr"`
	RequestURI       string           `json:"request_uri"`
	TransferEncoding []string         `json:"transfer_encoding"`
	ContentLength    int64            `json:"content_lenght"`
	Headers          PairResponseList `json:"headers"`
	Query            PairResponseList `json:"query"`
	FormData         PairResponseList `json:"form_data"`
	Body             string           `json:"body"`
}

func NewRecordResponse(record Record) RecordResponse {
	headers := make(PairResponseList, 0, len(record.recordData.headers))
	for _, p := range record.recordData.headers {
		pairResponse := PairResponse{
			Key:   p.key,
			Value: p.value,
		}
		headers = append(headers, pairResponse)
	}
	query := make(PairResponseList, 0, len(record.recordData.query))
	for _, p := range record.recordData.query {
		pairResponse := PairResponse{
			Key:   p.key,
			Value: p.value,
		}
		query = append(query, pairResponse)
	}
	recordResponse := RecordResponse{
		BinKey:           record.binKey.String(),
		RecordKey:        record.recordKey.String(),
		CreatedAt:        record.createdAt.Format(time.RFC3339),
		Method:           record.recordData.method,
		URL:              record.recordData.uRL,
		Proto:            record.recordData.proto,
		RemoteAddr:       record.recordData.remoteAddr,
		TransferEncoding: record.recordData.transferEncoding,
		ContentLength:    record.recordData.contentLength,
		Headers:          headers,
		Query:            query,
		FormData:         nil,
		Body:             string(record.recordData.body),
	}
	return recordResponse
}
