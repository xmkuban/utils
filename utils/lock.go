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
)

func init() {
	lockCache, _ = cache.NewCache("memory", `{"interval":30}`)
}

func GetServiceLock(service string, key interface{}) *sync.Mutex {
	k := service + fmt.Sprint(key)
	lock := lockCache.Get(k)
	if lock == nil {
		registerServiceLock.Lock()
		lock = lockCache.Get(k)
		if lock == nil {
			lock = new(sync.Mutex)
			lockCache.Put(k, lock, time.Second*30)
			registerServiceLock.Unlock()
			return lock.(*sync.Mutex)
		}
		registerServiceLock.Unlock()
	}
	return lock.(*sync.Mutex)
}
