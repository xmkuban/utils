package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xmkuban/utils/cache"
)

var (
	registerServiceLock = sync.Mutex{}
	lockCache           *cache.MemoryCache

	lockTimeout = 30 //server下的key锁过期时间, 单位秒, 目前限制最少存1s
)

func init() {
	lockCache = cache.NewMemoryCache()
}

func GetServiceLock(service string, key interface{}) *sync.Mutex {
	_lockTimeout := time.Second * time.Duration(lockTimeout)
	k := service + fmt.Sprint(key)
	lock, expire := lockCache.GetAndTTL(k)
	if lock == nil {
		registerServiceLock.Lock()
		lock = lockCache.Get(k)
		if lock == nil {
			lock = new(sync.Mutex)
			lockCache.Put(k, lock, _lockTimeout)
			registerServiceLock.Unlock()
			return lock.(*sync.Mutex)
		}
		registerServiceLock.Unlock()
	} else {
		if expire.Seconds() <= 1 {
			registerServiceLock.Lock()
			lockCache.Put(k, lock, _lockTimeout)
			registerServiceLock.Unlock()
		}
	}
	return lock.(*sync.Mutex)
}

// LockExtendTime 延长锁的使用，针对大时间事务需要延长锁的情况
func LockExtendTime(service string, key interface{}, lock *sync.Mutex) error {
	_lockTimeout := time.Second * time.Duration(lockTimeout)
	k := service + fmt.Sprint(key)
	_lock := lockCache.Get(k)
	if _lock == nil {
		return errors.New("not expired")
	}
	registerServiceLock.Lock()
	lockCache.Put(k, lock, _lockTimeout)
	registerServiceLock.Unlock()
	return nil
}

func SetKeyLockTimeout(t int) {
	if t <= 1 {
		return
	}
	lockTimeout = t
}
