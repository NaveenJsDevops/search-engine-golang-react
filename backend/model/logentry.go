package model

type LogEntry struct {
	MsgId          string `parquet:"name=MsgId"`
	PartitionId    int64  `parquet:"name=PartitionId"`
	Timestamp      string `parquet:"name=Timestamp"`
	Hostname       string `parquet:"name=Hostname"`
	Priority       int32  `parquet:"name=Priority"`
	Facility       int32  `parquet:"name=Facility"`
	FacilityString string `parquet:"name=FacilityString"`
	Severity       int32  `parquet:"name=Severity"`
	SeverityString string `parquet:"name=SeverityString"`
	AppName        string `parquet:"name=AppName"`
	ProcId         string `parquet:"name=ProcId"`
	Message        string `parquet:"name=Message"`
	MessageRaw     string `parquet:"name=MessageRaw"`
	StructuredData string `parquet:"name=StructuredData"`
	Tag            string `parquet:"name=Tag"`
	Sender         string `parquet:"name=Sender"`
	Groupings      string `parquet:"name=Groupings"`
	Event          string `parquet:"name=Event"`
	EventId        string `parquet:"name=EventId"`
	NanoTimeStamp  string `parquet:"name=NanoTimeStamp, type=UTF8"`
	Namespace      string `parquet:"name=namespace"`
}

type LogEntryRes struct {
	Limit      int64       `json:"limit"`
	OffSet     int64       `json:"offset"`
	Duration   string      `json:"duration"`
	LogEntries *[]LogEntry `json:"log"`
}
type LogEntryRecords struct {
	FileName         string `json:"fileName"`
	FileTotalRecords int64  `json:"fileRecords"`
}
type TotalLogs struct {
	TotalRecords    int64              `json:"totalRecords"`
	LogEntryRecords *[]LogEntryRecords `json:"logEntries"`
}
type Message struct {
	Message string `json:"messge"`
}

type LogEntryResAll struct {
	Duration   string      `json:"duration"`
	Total      int         `json:"total"`
	LogEntries *[]LogEntry `json:"log"`
}
