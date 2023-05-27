package background

import (
	"errors"
	c "health/clog"
	"sync"
	"time"
)

type RegisterInfo struct {
	tokenString string
	createdAt   int64
}

type RegisterMap struct {
	m map[string]*RegisterInfo
	l sync.Mutex
}

var regMap *RegisterMap

func RunRegMap() {

	regMap = &RegisterMap{m: make(map[string]*RegisterInfo)}

	go func() {
		var now time.Time
		ticker := time.NewTicker(time.Minute * 1)
		for ; true; now = <-ticker.C {
			if now.IsZero() {
				now = time.Now()
			}
			regMap.l.Lock()
			for k, v := range regMap.m {
				if now.Unix()-v.createdAt > 60 {
					c.InfoLog.Printf("backround: regmap delete expired key: %s", k)
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
		ri := &RegisterInfo{tokenString: tokenString, createdAt: time.Now().Unix()}
		regMap.m[key] = ri
		b = true
	} else {
		b = false
	}
	c.InfoLog.Printf("backround: regmap add register key: %s", key)
	regMap.l.Unlock()
	return
}

func (regMap *RegisterMap) Update(key string, tokenString string) (b bool) {
	regMap.l.Lock()
	regMap.m[key].tokenString = tokenString
	regMap.m[key].createdAt = time.Now().Unix()
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
