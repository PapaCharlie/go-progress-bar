package progressbar

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Progress struct {
	count      uint64
	lastReport time.Time

	Interval uint64
	MaxValue uint64

	lock *sync.Mutex
}

func NewProgressBar(interval, maxValue uint64) *Progress {
	if interval <= 0 {
		interval = 1 << 10
	}

	return &Progress{
		lastReport: time.Now(),
		lock:       &sync.Mutex{},
		Interval:   interval,
		MaxValue:   maxValue,
	}
}

var MAGS = []string{"", "K", "M", "G", "T", "P"}

func (p *Progress) Count() uint64 {
	return atomic.LoadUint64(&p.count)
}

func (p *Progress) Inc() bool {
	c := atomic.AddUint64(&p.count, 1)
	if c%p.Interval == 0 || c == p.MaxValue {
		p.lock.Lock()
		diff := time.Since(p.lastReport)
		p.lastReport = time.Now()
		speed := float64(p.Interval) / diff.Seconds()
		scaledSpeed := speed
		var mag string
		for _, h := range MAGS {
			if int64(scaledSpeed)/1000 > 0 {
				scaledSpeed = scaledSpeed / 1000
			} else {
				mag = h
				break
			}
		}
		if p.MaxValue > 0 {
			progress := float64(c) / float64(p.MaxValue) * 100
			eta := time.Duration(float64(p.MaxValue-c)/speed) * time.Second
			hours := int(eta.Hours())
			minutes := int(math.Mod(eta.Minutes(), 60))
			seconds := int(math.Mod(eta.Seconds(), 60))
			fmt.Printf("\rProgress: %7.3f%%, Speed: %5.1f %1sop/s, ETA: %02d:%02d:%02d",
				progress, scaledSpeed, mag, hours, minutes, seconds)
		} else {
			fmt.Printf("\rSpeed: %5.1f %sop/s", scaledSpeed, mag)

		}
		p.lock.Unlock()
		return true
	}
	return false
}
