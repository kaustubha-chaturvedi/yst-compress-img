package compressor

import (
	"fmt"
	"os"
	"sync"
	"time"
)

/* Progress Bar:
- We track: total images, completed images, and start time.
- Percentage = completed / total.
- ETA = based on average processing speed so far.
- Bar updates in-place using carriage return.
*/

type Progress struct {
	mu        sync.Mutex
	total     int
	completed int
	start     time.Time
}

func NewProgress(total int) *Progress {
	return &Progress{
		total: total,
		start: time.Now(),
	}
}

// mark image as done
func (p *Progress) Update() {
	p.mu.Lock()
	p.completed++
	p.print()
	p.mu.Unlock()
}

func (p *Progress) print() {
	percent := float64(p.completed) / float64(p.total)
	barLen := 20

	filled := int(percent * float64(barLen))
	empty := barLen - filled

	bar := "[" + recalc("#", filled) + recalc("-", empty) + "]"

	elapsed := time.Since(p.start).Seconds()
	imagesPerSec := float64(p.completed) / elapsed

	eta := ""
	if p.completed > 0 && imagesPerSec > 0 {
		remaining := float64(p.total-p.completed) / imagesPerSec
		eta = fmt.Sprintf(" | ETA %.1fs", remaining)
	}

	fmt.Fprintf(os.Stderr,
		"\r%s %3.0f%% | %d/%d images%s",
		bar,
		percent*100,
		p.completed,
		p.total,
		eta,
	)

	if p.completed == p.total {
		fmt.Fprintln(os.Stderr)
	}
}

func recalc(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
