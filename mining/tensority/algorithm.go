package tensority

// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L./zhizhusimd/ -l:cSimdTszhizhu.o -lstdc++ -lgomp -lpthread
// #include "./zhizhusimd/cSimdTs.h"
import "C"

import (
	"unsafe"

	"github.com/golang/groupcache/lru"

	"github.com/bytom/crypto/sha3pool"
	"github.com/bytom/protocol/bc"
	"time"
	"sync/atomic"
	"log"
)

const maxAIHashCached = 64

var countSimd uint64 = 1
const segementNumSimd = 20000  //the same with cSimdTs2.cpp resultarr

var lockArraySimd = [20000]int32{0}

func algorithm(blockHeader, seed *bc.Hash) *bc.Hash {
	timeNow := time.Now()


	bhBytes := blockHeader.Bytes()
	sdBytes := seed.Bytes()

	// Get thearray pointer from the corresponding slice
	bhPtr := (*C.uint8_t)(unsafe.Pointer(&bhBytes[0]))
	seedPtr := (*C.uint8_t)(unsafe.Pointer(&sdBytes[0]))


	var resPtr *C.uint8_t

	//resPtr := C.SimdTs2(bhPtr, seedPtr)

	num := atomic.AddUint64(&countSimd,1)
	if num >=1000000 {
		atomic.CompareAndSwapUint64(&countSimd,num,0)
	}else if num >= 2000000{
		//Force to 0
		atomic.StoreUint64(&countSimd,0)
	}

	mode := num%segementNumSimd

	swapRes := atomic.CompareAndSwapInt32(&(lockArraySimd[mode]),0,1); //add lock use flag

	if swapRes {
		resPtr = C.SimdTsNoLockSimd(bhPtr, seedPtr,C.uint32_t(mode))

		res := bc.NewHash(*(*[32]byte)(unsafe.Pointer(resPtr)))

		atomic.CompareAndSwapInt32(&(lockArraySimd[mode]),1,0);  //remove lock flag

		//log.Printf("remove lock flag %v  num:%v",mode,num)
		
		log.Printf("nolock use time %v,mode:%v,num:%v",time.Since(timeNow),mode,num)
		return &res
	}else {
		//there are 20000 req is handling,may be some problem happen.
		//log.Printf("poolloglock flag  userd by other mode:%v num:%v",mode,num)

		resPtr = C.SimdTsWithLockSimd(bhPtr, seedPtr,C.uint32_t(mode))

		res := bc.NewHash(*(*[32]byte)(unsafe.Pointer(resPtr)))
		
		log.Printf("lock use time %v,mode:%v",time.Since(timeNow),mode)
		return &res
	}
}


func calcCacheKey(hash, seed *bc.Hash) *bc.Hash {
	var b32 [32]byte
	sha3pool.Sum256(b32[:], append(hash.Bytes(), seed.Bytes()...))
	key := bc.NewHash(b32)
	return &key
}

// Cache is create for cache the tensority result
type Cache struct {
	lruCache *lru.Cache
}

// NewCache create a cache struct
func NewCache() *Cache {
	return &Cache{lruCache: lru.New(maxAIHashCached)}
}

// AddCache is used for add tensority calculate result
func (a *Cache) AddCache(hash, seed, result *bc.Hash) {
	key := calcCacheKey(hash, seed)
	a.lruCache.Add(*key, result)
}

// Hash is the real entry for call tensority algorithm
func (a *Cache) Hash(hash, seed *bc.Hash) *bc.Hash {
	key := calcCacheKey(hash, seed)
	if v, ok := a.lruCache.Get(*key); ok {
		return v.(*bc.Hash)
	}
	return algorithm(hash, seed)
}

// AIHash is created for let different package share same cache
var AIHash = NewCache()
