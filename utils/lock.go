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

	serverLockTimeout = time.Minute      //server锁过期时间
	keyLockTimeout    = time.Second * 30 //server下的key锁过期时间
)

func init() {
	lockCache, _ = cache.NewCache("memory", `{"interval":30}`)
}

func GetServiceLock(service string, key interface{}) *sync.Mutex {
	lockServer := lockCache.Get(service)
	if lockServer == nil {
		registerServiceLock.Lock()
		lockServer = lockCache.Get(service)
		if lockServer == nil {
			lockServer = new(sync.Mutex)
			lockCache.Put(service, lockServer, serverLockTimeout)
		}
		registerServiceLock.Unlock()
	}
	_lockServer := lockServer.(*sync.Mutex)
	k := service + fmt.Sprint(key)
	lock := lockCache.Get(k)
	if lock == nil {
		_lockServer.Lock()
		lock = lockCache.Get(k)
		if lock == nil {
			lock = new(sync.Mutex)
			lockCache.Put(k, lock, keyLockTimeout)
			_lockServer.Unlock()
			return lock.(*sync.Mutex)
		}
		_lockServer.Unlock()
	} else {
		_lockServer.Lock()
		lockCache.Put(k, lock, keyLockTimeout)
		_lockServer.Unlock()
	}
	return lock.(*sync.Mutex)
}

func SetServerLockTimeout(t time.Duration) {
	if t.Nanoseconds() == 0 {
		return
	}
	serverLockTimeout = t
}

func SetKeyLockTimeout(t time.Duration) {
	if t.Nanoseconds() == 0 {
		return
	}
	keyLockTimeout = t
}
