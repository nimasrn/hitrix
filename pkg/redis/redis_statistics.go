package redis

import (
	"strings"

	"github.com/latolukasz/fluxaorm"
)

type Statistics struct {
	RedisPool string
	Info      map[string]string
}

func GetRedisStatistics(ormService fluxaorm.Context, dragonflyDBPools map[string]struct{}) []*Statistics {
	pools := ormService.Engine().Registry().RedisPools()
	results := make([]*Statistics, 0)

	for pool, redis := range pools {
		infoSection := "everything"

		_, has := dragonflyDBPools[pool]
		if has {
			infoSection = "all"
		}

		poolStats := &Statistics{RedisPool: pool, Info: make(map[string]string)}

		info := redis.Info(ormService, infoSection)
		lines := strings.Split(info, "\r\n")

		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			row := strings.Split(line, ":")
			val := ""

			if len(row) > 1 {
				val = row[1]
			}

			poolStats.Info[row[0]] = val
		}

		results = append(results, poolStats)
	}

	return results
}
