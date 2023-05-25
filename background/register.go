package background

import (
	"errors"
	"sync"
	"time"
)

type RegisterInfo struct {
	tokenString string
	expiry      int64
}

type RegisterMap struct {
	m map[string]*RegisterInfo
	l sync.Mutex
}

var regMap *RegisterMap

func RunRegMap() {

	regMap = &RegisterMap{m: make(map[string]*RegisterInfo)}

	go func() {
		for now := range time.Tick(time.Minute) {
			regMap.l.Lock()
			for k, v := range regMap.m {
				if v.expiry-now.Unix() < 0 {
					delete(regMap.m, k)
				}
			}
			regMap.l.Unlock()
		}
	}()
}

func GetRegMap() *RegisterMap {
	return regMap
}

func (regMap *RegisterMap) Len() int {
	return len(regMap.m)
}

func (regMap *RegisterMap) Add(key string, tokenString string) (b bool) {
	regMap.l.Lock()
	if _, ok := regMap.m[key]; !ok {
		ri := &RegisterInfo{tokenString: tokenString, expiry: time.Now().Unix() + (180)}
		regMap.m[key] = ri
		b = true
	} else {
		b = false
	}
	regMap.l.Unlock()
	return
}

func (regMap *RegisterMap) Update(key string, tokenString string) (b bool) {
	regMap.l.Lock()
	regMap.m[key].tokenString = tokenString
	regMap.m[key].expiry = time.Now().Unix() + (120)
	regMap.l.Unlock()
	return
}

func (regMap *RegisterMap) Get(key string) (v *RegisterInfo, e error) {
	regMap.l.Lock()
	if ri, ok := regMap.m[key]; ok {
		v = ri
		e = nil
	} else {
		v = nil
		e = errors.New("key not found")
	}
	regMap.l.Unlock()
	return
}

func (regMap *RegisterMap) Exists(key string) (b bool) {
	regMap.l.Lock()
	if _, ok := regMap.m[key]; ok {
		b = true
	} else {
		b = false
	}
	regMap.l.Unlock()
	return
}

func (regMap *RegisterMap) Delete(key string) {
	regMap.l.Lock()
	delete(regMap.m, key)
	regMap.l.Unlock()
}
