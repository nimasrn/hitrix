package entity

import (
	"github.com/latolukasz/fluxaorm"
)

type OSSBucketCounterEntity struct {
	ID      uint64 `orm:"table=oss_buckets_counters"`
	Counter uint64 `orm:"required"`
}
