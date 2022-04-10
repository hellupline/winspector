package responses

import (
	"time"

	"github.com/hellupline/winspector/pkg/models"
)

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
	PostFormData     PairResponseList `json:"post_form_data"`
	Body             string           `json:"body"`
}

func NewRecordResponse(record models.Record) RecordResponse {
	headers := make(PairResponseList, 0, len(record.RequestData.Headers))
	for _, p := range record.RequestData.Headers {
		pairResponse := PairResponse{Key: p.Key, Value: p.Value}
		headers = append(headers, pairResponse)
	}
	query := make(PairResponseList, 0, len(record.RequestData.Query))
	for _, p := range record.RequestData.Query {
		pairResponse := PairResponse{Key: p.Key, Value: p.Value}
		query = append(query, pairResponse)
	}
	postFormData := make(PairResponseList, 0, len(record.RequestData.PostFormData))
	for _, p := range record.RequestData.PostFormData {
		pairResponse := PairResponse{Key: p.Key, Value: p.Value}
		postFormData = append(postFormData, pairResponse)
	}
	recordResponse := RecordResponse{
		BinKey:           record.BinKey.String(),
		RecordKey:        record.RecordKey.String(),
		CreatedAt:        record.CreatedAt.Format(time.RFC3339),
		Method:           record.RequestData.Method,
		URL:              record.RequestData.URL,
		Proto:            record.RequestData.Proto,
		Host:             record.RequestData.Host,
		RemoteAddr:       record.RequestData.RemoteAddr,
		RequestURI:       record.RequestData.RequestURI,
		TransferEncoding: record.RequestData.TransferEncoding,
		ContentLength:    record.RequestData.ContentLength,
		Headers:          headers,
		Query:            query,
		PostFormData:     postFormData,
		Body:             string(record.RequestData.Body),
	}
	return recordResponse
}
