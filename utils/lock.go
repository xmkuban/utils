package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/xmkuban/utils/cache"
)

var (
	registerServiceLock = sync.Mutex{}
	lockCache           cache.Cache

	lockTimeout = 30 //server下的key锁过期时间, 单位秒, 目前限制最少存1s
)

func init() {
	lockCache, _ = cache.NewCache("memory", `{"interval":30}`)
}

func GetServiceLock(service string, key interface{}) *sync.Mutex {
	_lockTimeout := time.Second * time.Duration(lockTimeout)
	k := service + fmt.Sprint(key)
	lock := lockCache.Get(k)
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
		expire := lockCache.TTL(k)
		if expire.Seconds() <= 1 {
			registerServiceLock.Lock()
			lockCache.Put(k, lock, _lockTimeout)
			registerServiceLock.Unlock()
		}
	}
	return lock.(*sync.Mutex)
}

func SetKeyLockTimeout(t int) {
	if t <= 1 {
		return
	}
	lockTimeout = t
}
