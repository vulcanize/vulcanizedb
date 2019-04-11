package memsize

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBitmapBlock(t *testing.T) {
	marks := map[uintptr]bool{
		10:  true,
		13:  true,
		44:  true,
		128: true,
		129: true,
		256: true,
		700: true,
	}
	var b bmBlock
	for i := range marks {
		b.mark(i)
	}
	for i := uintptr(0); i < bmBlockRange; i++ {
		if b.isMarked(i) && !marks[i] {
			t.Fatalf("wrong mark at %d", i)
		}
	}
	if count := b.count(0, bmBlockRange-1); count != len(marks) {
		t.Fatalf("wrong onesCount: got %d, want %d", count, len(marks))
	}
}

func TestBitmapBlockCount(t *testing.T) {
	var b bmBlock
	// Mark addresses (90,250)
	for i := 90; i < 250; i++ {
		b.mark(uintptr(i))
	}
	// Check counts.
	tests := []struct {
		start, end uintptr
		want       int
	}{
		{start: 0, end: 0, want: 0},
		{start: 0, end: 10, want: 0},
		{start: 0, end: 250, want: 160},
		{start: 0, end: 240, want: 150},
		{start: 0, end: bmBlockRange - 1, want: 160},
		{start: 100, end: bmBlockRange - 1, want: 150},
		{start: 100, end: 110, want: 10},
		{start: 100, end: 250, want: 150},
		{start: 100, end: 211, want: 111},
		{start: 111, end: 211, want: 100},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-%d", test.start, test.end), func(t *testing.T) {
			if c := b.count(test.start, test.end); c != test.want {
				t.Errorf("wrong onesCountRange(%d, %d): got %d, want %d", test.start, test.end, c, test.want)
			}
		})
	}
}

func TestBitmapMarkRange(t *testing.T) {
	N := 1000

	// Generate random non-overlapping mark ranges.
	var (
		r      = rand.New(rand.NewSource(312321312))
		bm     = newBitmap()
		ranges = make(map[uintptr]uintptr)
		addr   uintptr
		total  uintptr // number of bytes marked
	)
	for i := 0; i < N; i++ {
		addr += uintptr(r.Intn(bmBlockRange))
		len := uintptr(r.Intn(40))
		total += len
		ranges[addr] = len
		bm.markRange(addr, len)
	}

	// Check all marks are set.
	for start, len := range ranges {
		for i := uintptr(0); i < len; i++ {
			if !bm.isMarked(start + i) {
				t.Fatalf("not marked at %d", start)
			}
		}
	}

	// Check total number of bits is reported correctly.
	if c := bm.countRange(0, addr+ranges[addr]); c != total {
		t.Errorf("countRange(0, %d) returned %d, want %d", addr, c, total)
	}

	// Probe random addresses.
	for i := 0; i < N; i++ {
		addr := uintptr(r.Uint64())
		marked := false
		for start, len := range ranges {
			if addr >= start && addr < start+len {
				marked = true
				break
			}
		}
		if bm.isMarked(addr) && !marked {
			t.Fatalf("extra mark at %d", addr)
		}
	}
}

func BenchmarkBitmapMarkRange(b *testing.B) {
	var addrs [2048]uintptr
	r := rand.New(rand.NewSource(423098209802))
	for i := range addrs {
		addrs[i] = uintptr(r.Uint64())
	}

	doit := func(b *testing.B, rlen int) {
		bm := newBitmap()
		for i := 0; i < b.N; i++ {
			addr := addrs[i%len(addrs)]
			bm.markRange(addr, uintptr(rlen))
		}
	}
	for rlen := 1; rlen <= 4096; rlen *= 8 {
		b.Run(fmt.Sprintf("%d", rlen), func(b *testing.B) { doit(b, rlen) })
	}
}
