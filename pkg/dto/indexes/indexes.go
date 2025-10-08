package indexes

type ResponseDTOList struct {
	Indexes []Index
}

type Index struct {
	Name      string
	TotalDocs int
	TotalSize uint64
}
