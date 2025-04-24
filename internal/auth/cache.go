package auth

import "time"

func UserCacheCleaner() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for _ = range ticker.C {
			UserCachesMu.Lock()
			for k, v := range UserCaches {
				if v.Expired() {
					delete(UserCaches, k)
				}
			}
			UserCachesMu.Unlock()
		}
	}()
}
