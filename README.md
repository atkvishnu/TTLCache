# TTLCache
### In-memory cache with per-key expiration in Golang.    

> Map is used for the cache.    
> Mutex is used to synchronize concurrent access.    
> Goroutine used to check cache expiration.

### Methods

- New() - Create a new cache.
- Get(key) - Get a value from the cache based on the key.
- Set(Key, Value, ExpirationTime) - Set key-value pairs with TTL.

---

### Run

$ `go run TTLCache.go`


---

### TODO: 
- [ ] - Make test files.  
