package delayedqueue

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xorcare/pointer"

	"github.com/coretrix/hitrix/pkg/dto/delayedqueue"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context) *delayedqueue.List {
	appService := service.DI().App()
	ormService := service.DI().OrmForContext(ctx)

	result := &delayedqueue.List{Rows: make([]delayedqueue.Row, len(appService.RedisDelayedQueues))}
	for i, queue := range appService.RedisDelayedQueues {
		result.Rows[i].Queue = queue
		result.Rows[i].Total = ormService.Engine().Redis(appService.RedisPools.Stream).ZCount(ormService, queue, "-inf", "+inf")
		values := ormService.Engine().Redis(appService.RedisPools.Stream).ZRangeArgsWithScores(ormService, redis.ZRangeArgs{
			Key:     queue,
			Start:   0,
			Stop:    "+inf",
			Offset:  0,
			Count:   1,
			Rev:     true,
			ByScore: true})

		if len(values) > 0 {
			result.Rows[i].LatestItem = pointer.Uint64(uint64(values[0].Score))
		}
	}

	return result
}
