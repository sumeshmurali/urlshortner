package database

type UrlDetail struct {
	ID         int
	Url        string
	LongUrl    string
	VisitCount int
}

type UrlMeta struct {
	UrlId       int
	Ip          string
	Location    string
	Device_type string
}
