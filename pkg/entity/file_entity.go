package entity

import (
	"time"
)

type FileStatus string

func (f FileStatus) String() string {
	return string(f)
}

const (
	FileStatusNew       FileStatus = "new"
	FileStatusProcessed FileStatus = "processed"
)

type fileStatus struct {
	FileStatusNew       string
	FileStatusProcessed string
}

var FileStatusAll = fileStatus{
	FileStatusNew:       FileStatusNew.String(),
	FileStatusProcessed: FileStatusProcessed.String(),
}

type FileObject struct {
	ID         uint64
	StorageKey string
	Data       interface{}
}

type FileEntity struct {
	ID        uint64 `orm:"table=files;redisSearch=search_pool;searchable;sortable"`
	File      *FileObject
	Status    string    `orm:"required;enum=entity.FileStatusAll"`
	Namespace string    `orm:"required;searchable"`
	CreatedAt time.Time `orm:"time=true"`
}
