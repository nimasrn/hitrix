package locker

import (
	"errors"
	"time"

	"github.com/latolukasz/fluxaorm"
	tusd "github.com/tus/tusd/pkg/handler"

	"github.com/coretrix/hitrix/service"
)

type RedisLocker struct {
	ormService fluxaorm.Context
}

func (locker *RedisLocker) NewLock(id string) (tusd.Lock, error) {
	return &redisLock{id: id, ormService: locker.ormService}, nil
}

type redisLock struct {
	id         string
	ormService fluxaorm.Context
	redisLock  *fluxaorm.Lock
}

func (lock *redisLock) Lock() error {
	redisLock, obtained := lock.ormService.Engine().Redis(service.DI().App().RedisPools.Cache).GetLocker().Obtain(
		lock.ormService,
		"tusd:upload:lock:"+lock.id, time.Hour*24, time.Second*2,
	)
	if !obtained {
		return errors.New("cannot obtain lock")
	}

	lock.redisLock = redisLock

	return nil
}

func (lock *redisLock) Unlock() error {
	lock.redisLock.Release(lock.ormService)

	return nil
}
