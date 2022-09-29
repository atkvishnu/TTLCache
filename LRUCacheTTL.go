/**
APPROACH:
- Implemented a basic LRU cache using a hashmap and a doubly linked list.
- Key in map is the key, the Value in map is the address to the Node of the Doubly Linked List.
- Nodes of the Doubly Linked List have values such as Prev*, Key, Value, TTL, Next*
- Implemented a ADT doubly linked list with some methods like:  AddAtBeg(k, v, ttl), AddAtEnd(k, v, ttl), whatsAt(k), DelAtBeg(), DelNode(k), DelLRUNode(), Count(), Display(). Did not use the `container/list` library.
- The LRUCache functions works like this:
	- Set(key, value, ttl): + Will check if key is present in hashmap.
				+ If Key is present, do not add anything and return "key already present in map"
				+ If Key not present, then check if hashmap is empty for insertion (size check - edge case).
				+ If Hashmap is empty, Push the Node in front of Doubly Linked List, and add the key and address of Node(as value) in hashmap.
				+ If Hashmap is not empty, then remove the Least Recently Used (LRU) element from the back of the Doubly Linked List, and then add the new Node to the front of DLL and HashMap.
				+ We do not have a eviction function as these 2 functions will help us comply with the LRU eviction policy.
				+ This completes the 1st TODO ✔️
	- Get(key) : + Will check if the key is present in HashMap.
		     + If Key is not present, return "not found"
		     + If Key is found, then dereference the address which is stored in the value section of HashMap to obtain the values inside the Nodes, as we touched this Node therefore we need to delete it from it's place and put it in front of the DLL again, to comply with the LRU eviction policy.
		     + Return the Value associated with key.
       		     + We don't need to check if the Key is expired or not because our Goroutine will be running every 1 Seconds which will take care of the expired Nodes.
                     + This completes our 2nd TODO ✔
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

/*---------------------------Doubly-Linked-List-Implementation-Starts------------------------------*/

// Node struct
type Node struct {
	Key  any
	Val  any
	TTL  int64
	Prev *Node
	Next *Node
}

// NewNode Creates a new node.
func NewNode(k1 any, v1 any, TTL1 int64) *Node {
	return &Node{k1, v1, TTL1, nil, nil}
}

type Doubly struct {
	Head *Node
}

func NewDoubly() *Doubly {
	return &Doubly{nil}
}

// AddAtBeg adds a node to the beginning of the doubly linked list
func (ll *Doubly) AddAtBeg(k any, v any, ttl int64) **Node {
	n := NewNode(k, v, ttl)
	n.Next = ll.Head

	if ll.Head != nil {
		ll.Head.Prev = n
	}

	ll.Head = n

	return &ll.Head
}

// AddAtEnd adds a node at the end of the doubly linked list
func (ll *Doubly) AddAtEnd(k any, v any, ttl int64) {
	n := NewNode(k, v, ttl)

	if ll.Head == nil {
		ll.Head = n
		return
	}

	cur := ll.Head
	for ; cur.Next != nil; cur = cur.Next {
	}
	cur.Next = n
	n.Prev = cur
}

// AddAtEnd adds a node at the end of the doubly linked list
func (ll *Doubly) whatsAt(k any) *Node {
	cur := ll.Head
	for ; cur.Next != nil; cur = cur.Next {
		if cur.Next.Key == k {
			it := NewNode(cur.Next.Key, cur.Next.Val, cur.Next.TTL)
			return it
		}
	}

	return nil
}

// DelAtBeg deletes the node at the beginning of the doubly linked list
func (ll *Doubly) DelAtBeg() any {
	if ll.Head == nil {
		return -1
	}

	cur := ll.Head
	ll.Head = cur.Next

	if ll.Head != nil {
		ll.Head.Prev = nil
	}
	return cur.Val
}

// DeleteNode finds the given key and deleted that node, if it's present
func (ll *Doubly) DeleteNode(k any) *Node {
	if ll.Head == nil {
		return nil
	}
	// only one item
	if ll.Head.Next == nil {
		ll.DelAtBeg()
	}

	// more than one, go to second last
	cur := ll.Head
	for ; cur.Next.Next != nil; cur = cur.Next {
		if cur.Next.Key == k && cur.Next != nil {
			retVal := NewNode(cur.Next.Key, cur.Next.Val, cur.Next.TTL)
			cur.Next.Next.Prev = cur
			cur.Next = cur.Next.Next
			return retVal
		}
	}
	return nil
}

// DelLRUNode Delete a node at the end of the doubly linked list
// DelLRUNode later
func (ll *Doubly) DelLRUNode() *Node {
	// no item
	if ll.Head == nil {
		return nil
	}

	// only one item
	if ll.Head.Next == nil {
		ll.DelAtBeg()
	}

	// more than one, go to second last
	cur := ll.Head
	for ; cur.Next.Next != nil; cur = cur.Next {
	}

	returnVal := NewNode(cur.Next.Key, cur.Next.Val, cur.Next.TTL)
	cur.Next = nil

	return returnVal
}

// Count Number of nodes in the doubly linked list
func (ll *Doubly) Count() any {
	var ctr = 0

	for cur := ll.Head; cur != nil; cur = cur.Next {
		ctr += 1
	}

	return ctr
}

// Display the doubly linked list
func (ll *Doubly) Display() {
	for cur := ll.Head; cur != nil; cur = cur.Next {
		fmt.Print(cur.Key, " ", cur.Val, " ", cur.TTL)
	}
	fmt.Println("")
}

/*---------------------------Doubly-Linked-List-Implementation-Ends------------------------------*/

/*----------------------LRU-Cache-with-TTL-feature-implementation-starts-------------------------*/

// Cache is a map that auto-expires
type Cache struct {
	cacheMap   map[any]any
	helperList *Doubly
	mutex      sync.RWMutex
}

// NewCache is a helper to create instance of Cache
func NewCache(size int) *Cache {
	cache := &Cache{
		cacheMap:   make(map[any]any, size),
		helperList: NewDoubly(),
		mutex:      sync.RWMutex{},
	}
	return cache
}

// LenMap returns the number of cacheMap in the cache
func (cache *Cache) LenMap() any {
	cache.mutex.Lock()
	var Len = len(cache.cacheMap)
	cache.mutex.Unlock()
	return Len
}

// Set is a way to add new cacheMap to the map
func (cache *Cache) Set(k any, v any, ttl int64, maxSize int) any {

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	_, present := cache.cacheMap[k]
	if present { // if key already present in map
		return "Key is already present!"
	} else {
		if len(cache.cacheMap) < maxSize { // if hashmap is empty, then add in front of DLL
			adr := cache.helperList.AddAtBeg(k, v, ttl) // add node in front of DLL
			cache.cacheMap[k] = adr                     // add address of the new Node as value in hashmap
			return "Node Added!"
		} else { // if hashmap is full, then 	// eviction policy
			deletedNode := cache.helperList.DelLRUNode()    // remove LRU (from back of DLL)
			delete(cache.cacheMap, deletedNode.Key)     // and delete K-V pair from hashmap
			adr := cache.helperList.AddAtBeg(k, v, ttl) // add node in front of DLL
			cache.cacheMap[k] = adr                     // add address of the new Node as value in hashmap
			return "Node Added after removing one node!"
		}
	}
}

// Get method to lookup cacheMap (Every lookup, also increases TTL of the item, hence extending its life)
func (cache *Cache) Get(key any) (value *Node, found bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	_, exists := cache.cacheMap[key] // check if key exists in hashmap
	if !exists {
		value = nil
		found = false
	} else { // if key found
		nodeValue := cache.helperList.DeleteNode(key) // Delete that node and
		delete(cache.cacheMap, nodeValue.Key)
		newTime := time.Now().Unix()                                                // we are touching the file, therefore when we add it again to the DLL we are refreshing the time
		newAddr := cache.helperList.AddAtBeg(nodeValue.Key, nodeValue.Val, newTime) // Add it at the beginning of the DLL
		cache.cacheMap[nodeValue.Key] = newAddr
		return nodeValue, true
	}
	return
}

func main() {
	maxSize := 1024
	timeToLife := 7

	cache := NewCache(maxSize)

	// Setting Random Key-Value pairs
	for i := 0; i < 1024; i++ {
		k := i
		v := i * 75
		t := time.Now().Unix()
		cache.Set(k, v, t, maxSize)
	}

	// fn `time.NewTicker()` makes a channel that sends a periodic message, and provides a way to stop it.
    


    // goroutine
	go func() {
		for now := range time.Tick(time.Second * 1) { // ticker to check every second (can manipulate the cleaning timer)
			cache.mutex.Lock()
			ll := cache.helperList
			for cur := ll.Head; cur.Next != nil; cur = cur.Next {
				if now.Unix()-cur.TTL > int64(timeToLife) {
					cache.helperList.DeleteNode(cur.Key) // delete from DLL
					delete(cache.cacheMap, cur.Key)      // delete from map
				}
			}
			cache.mutex.Unlock()
		}
	}()

	// wg.Add(10)

	fmt.Println("")

	fmt.Println("Initial length of Cache (len(m)):", cache.LenMap())
	fmt.Println("")

	fmt.Println("Elements of map: ", cache.cacheMap)
	fmt.Println("")

	fmt.Println("-----5 second sleep starts----- ")
	time.Sleep(5 * time.Second)
	fmt.Println("Length of Cache after 5 seconds (this won't be empty):", cache.LenMap())
	fmt.Println("")
	// fmt.Println("Elements in map: ", cache.cacheMap)

	// Setting Random Key-Value pairs
	for i := 50; i < 100; i++ {
		val, bl := cache.Get(i)
		if !bl {
			fmt.Println("Key not found: ")
		} else {
			fmt.Println("Key Found in map! Key: ", val.Key, ", Value: ", val.Val, " TTL: ", val.TTL)
		}
	}
	fmt.Println("\n Now we will be having 50 keys + 1 (Head) in map (50-99) with refreshed TTL! ")

	fmt.Println("-----3 second sleep starts----- ")
	time.Sleep(3 * time.Second)
	fmt.Println("Final length of Cache after 8 seconds (this should have 51 Nodes):", cache.LenMap())
	fmt.Println("Elements in map (empty map): ", cache.cacheMap)
	fmt.Println("")

	fmt.Println("-----2 second sleep starts----- ")
	time.Sleep(2 * time.Second)
	fmt.Println("Final length of Cache after 10 seconds (this should have 1 Node):", cache.LenMap())
	fmt.Println("Elements in map (empty map): ", cache.cacheMap)
	fmt.Println("")

	fmt.Println("-----9 second sleep starts----- ")
	time.Sleep(9 * time.Second)
	fmt.Println("Final length of Cache after 19 seconds (this will be 1, because Head is not deleted):", cache.LenMap())
	fmt.Println("Elements in map (empty map): ", cache.cacheMap)
	fmt.Println("")

	val, bl := cache.Get(1)
	if !bl {
		fmt.Println("Key not found")
	} else {
		fmt.Println("Key found in map, value related to key is: ", val)
	}

}
