package common

type AccessKey struct {
	Id     string
	Secret string
}

type RecordInfo struct {
	Domain string
	Name   string
	Id     string
}

type RecordHandler interface {
	SetRecordIp(info *RecordInfo, ip string) error
}
