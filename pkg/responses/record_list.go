package responses

import (
	"sort"

	"github.com/hellupline/winspector/pkg/models"
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

func NewRecordResponseList(records []models.Record) RecordResponseList {
	recordResponseList := make(RecordResponseList, 0, len(records))
	for _, record := range records {
		recordResponse := NewRecordResponse(record)
		recordResponseList = append(recordResponseList, recordResponse)
	}
	sort.Sort(recordResponseList)
	return recordResponseList
}
