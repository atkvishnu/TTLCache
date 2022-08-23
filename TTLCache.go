/**
 * @brief      TTL in-memory cache (LogicQ)
 * 				
 * @hints	   - Use a map for the cache and synchronize concurrent access with a synchronization primitive like mutex.
 *  			 Go offers a few mutex options and pick one that works best for performance.
 * 			   - Use a goroutine to handle and check cache expirations.
 * 			   
 * @methods	   - New() - Create a new cache
 *			   - Get(key) - Get a value from the cache based on the key. An expired key or non-existent key should return a "not found".
 * 			   - Set(Key, value, Expiration)
 * 			   
 * @author     Kumar Vishnu
 * @date       21-08-2022
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

// item struct
type item struct {
	value      int
	lastAccess int64
}

// TTLCacheMap is our custom struct storing a map of item(struct) and mutex.
type TTLCacheMap struct {
	cacheMap map[int]*item 		// map of item
	mutex sync.RWMutex  		// RWmutex for concurrent access
}


/**
 * @brief      `New` returns a new TTLCacheMap with the given size
 *
 * @param      size  The size of cache
 *
 * @return     m is a pointer to struct TTLCacheMap
 */
func New(size int) (m *TTLCacheMap) {
	m = &TTLCacheMap{cacheMap: make(map[int]*item, size)}
	return			// (naked return)
}


/**
 * @brief      Set(key, value, timeToLife) - will set the key, value, and expiration time in map if key is not available in the map 
 *
 * @param      key         The key
 * @param      value       The value
 * @param      timeToLife  Expiration time
 */
func (m *TTLCacheMap) Set(key int, value int, timeToLife int) {
	m.mutex.Lock()         
	it, ok := m.cacheMap[key] 			// check if key is in map (ok is bool) and assign to it (pointer to item) if it is in map or nil if not in map
	if !ok {               				// if not in map
		it = &item{value: value} 		// create new item with value
		m.cacheMap[key] = it        	// add to map
	}
	it.lastAccess = time.Now().Unix() 	// update last access time to current time in unix seconds
	m.mutex.Unlock() 
}


/**
 * @brief      Get(key) returns the value and updates the lastAccess time
 *
 * @param      key     The key
 * @param      output  The value associated to the key
 *
 * @return     Value associated to given key
 */
func (m *TTLCacheMap) Get(key int) (output int) {
	m.mutex.Lock()
	if it, ok := m.cacheMap[key]; ok { 	  	// if key is in map (ok is bool) and assign to it (pointer to item) if it is in map or nil if not in map
		output = it.value                 	// assign value to output
		it.lastAccess = time.Now().Unix() 	// update last access time to current time in unix seconds (int64)
	}
	m.mutex.Unlock()
	return
}


/**
 * @brief      Len returns the length of the map
 */
func (m *TTLCacheMap) Len() int {
	return len(m.cacheMap)
}


func main() {
	m := New(1024)  						// create new cache with size 1024 and assigned to m
	timeToLife := 7 						// time to live = 7 seconds

	// Setting Random Key-Value pairs
	for i := 0; i < 100; i++ {
		key := i
		value := i * 75
		m.Set(key, value, timeToLife)
	}
	fmt.Println("\n")
	
	fmt.Println("Initial length of Cache (len(m)):", m.Len())
	fmt.Println("\n")

	fmt.Println("Elements of map: ", m.cacheMap)
	fmt.Println("\n")

	// `goroutine` to delete expired items in map (seperate goroutine initiated)
	go func() {
		for now := range time.Tick(time.Second * 0.5) { 							// ticker to check every second (can manipulate the cleaning timer)
			m.mutex.Lock()                    
			for key, value := range m.cacheMap { 								// loop through map
				if now.Unix()-value.lastAccess > int64(timeToLife) { 			// if last access time is greater than timeToLive
					delete(m.cacheMap, key) 									// delete from map
				}
			}
			m.mutex.Unlock()
		}
	}()

	fmt.Println("-----5 second sleep starts----- \n")
	time.Sleep(5 * time.Second)
	fmt.Println("Length of Cache after 5 seconds (this won't be empty):", m.Len())
	fmt.Println("\n")
	// fmt.Println("Elements in map: ", m.cacheMap)

	fmt.Println("-----4 second sleep starts----- \n")
	time.Sleep(4 * time.Second)
	fmt.Println("Final length of Cache after 9 seconds (this will be empty):", m.Len())
	fmt.Println("Elements in map (empty map): ", m.cacheMap)
}



