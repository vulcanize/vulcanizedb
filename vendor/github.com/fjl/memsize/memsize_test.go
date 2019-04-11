package memsize

import (
	"testing"
	"unsafe"
)

const (
	sizeofSlice     = unsafe.Sizeof([]byte{})
	sizeofMap       = unsafe.Sizeof(map[string]string{})
	sizeofInterface = unsafe.Sizeof((interface{})(nil))
	sizeofString    = unsafe.Sizeof("")
	sizeofWord      = unsafe.Sizeof(uintptr(0))
	sizeofChan      = unsafe.Sizeof(make(chan struct{}))
)

type (
	struct16 struct {
		x, y uint64
	}
	structptr struct {
		x   uint32
		cld *structptr
	}
	structuint32ptr struct {
		x *uint32
	}
	structmultiptr struct {
		s1 *structptr
		u1 *structuint32ptr
		s2 *structptr
		u2 *structuint32ptr
		s3 *structptr
		u3 *structuint32ptr
	}
	structarrayptr struct {
		x *uint64
		a [10]uint64
	}
	structiface struct {
		s *struct16
		x interface{}
	}
	struct64array  struct{ array64 }
	structslice    struct{ s []uint32 }
	structstring   struct{ s string }
	structloop     struct{ s *structloop }
	structptrslice struct{ s *structslice }
	array64        [64]byte
)

func TestTotal(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want uintptr
	}{
		{
			name: "struct16",
			v:    &struct16{},
			want: 16,
		},
		{
			name: "structptr_nil",
			v:    &structptr{},
			want: 2 * sizeofWord,
		},
		{
			name: "structptr",
			v:    &structptr{cld: &structptr{}},
			want: 2 * 2 * sizeofWord,
		},
		{
			name: "structptr_loop",
			v: func() *structptr {
				v := &structptr{}
				v.cld = v
				return v
			}(),
			want: 2 * sizeofWord,
		},
		{
			name: "structmultiptr_loop",
			v: func() *structmultiptr {
				v1 := &structptr{x: 1}
				v2 := &structptr{x: 2, cld: v1}
				return &structmultiptr{s1: v1, s2: v1, s3: v2}
			}(),
			want: 6*sizeofWord /* structmultiptr */ + 2*2*sizeofWord, /* structptr */
		},
		{
			name: "structmultiptr_interior",
			v: func() *structmultiptr {
				v1 := &structptr{x: 1}
				v2 := &structptr{x: 2}
				return &structmultiptr{
					// s1 is scanned before u1, which has a reference to a field of s1.
					s1: v1,
					u1: &structuint32ptr{x: &v1.x},
					// This one goes the other way around: u2, which has a reference to a
					// field of s3 is scanned before s3.
					u2: &structuint32ptr{x: &v2.x},
					s3: v2,
				}
			}(),
			want: 6*sizeofWord /* structmultiptr */ + 2*2*sizeofWord /* structptr */ + 2*sizeofWord, /* structuint32ptr */
		},
		{
			name: "struct64array",
			v:    &struct64array{},
			want: 64,
		},
		{
			name: "structptrslice",
			v:    &structptrslice{&structslice{s: []uint32{1, 2, 3}}},
			want: sizeofWord + sizeofSlice + 3*4,
		},
		{
			name: "array_unadressable",
			v: func() *map[[3]uint64]struct{} {
				v := map[[3]uint64]struct{}{
					{1, 2, 3}: struct{}{},
				}
				return &v
			}(),
			want: sizeofMap + 3*8,
		},
		{
			name: "structslice",
			v:    &structslice{s: []uint32{1, 2, 3}},
			want: sizeofSlice + 3*4,
		},
		{
			name: "structloop",
			v: func() *structloop {
				v := new(structloop)
				v.s = v
				return v
			}(),
			want: sizeofWord,
		},
		{
			name: "array64",
			v:    &array64{},
			want: 64,
		},
		{
			name: "byteslice",
			v:    &[]byte{1, 2, 3},
			want: sizeofSlice + 3,
		},
		{
			name: "slice3_ptrval",
			v:    &[]*struct16{{}, {}, {}},
			want: sizeofSlice + 3*sizeofWord + 3*16,
		},
		{
			name: "map3",
			v:    &map[uint64]uint64{1: 1, 2: 2, 3: 3},
			want: sizeofMap + 3*8 /* keys */ + 3*8, /* values */
		},
		{
			name: "map3_ptrval",
			v:    &map[uint64]*struct16{1: {}, 2: {}, 3: {}},
			want: sizeofMap + 3*8 /* keys */ + 3*sizeofWord /* value pointers */ + 3*16, /* values */
		},
		{
			name: "map3_ptrkey",
			v:    &map[*struct16]uint64{{x: 1}: 1, {x: 2}: 2, {x: 3}: 3},
			want: sizeofMap + 3*sizeofWord /* key pointers */ + 3*16 /* keys */ + 3*8, /* values */
		},
		{
			name: "map_interface",
			v:    &map[interface{}]interface{}{"aa": uint64(1)},
			want: sizeofMap + sizeofInterface + sizeofString + 2 /* key */ + sizeofInterface + 8, /* value */
		},
		{
			name: "pointerpointer",
			v: func() **uint64 {
				i := uint64(0)
				p := &i
				return &p
			}(),
			want: sizeofWord + 8,
		},
		{
			name: "structstring",
			v:    &structstring{"123"},
			want: sizeofString + 3,
		},
		{
			name: "slices_samearray",
			v: func() *[3][]byte {
				backarray := [64]byte{}
				return &[3][]byte{
					backarray[16:],
					backarray[4:16],
					backarray[0:4],
				}
			}(),
			want: 3*sizeofSlice + 64,
		},
		{
			name: "slices_nil",
			v: func() *[2][]byte {
				return &[2][]byte{nil, nil}
			}(),
			want: 2 * sizeofSlice,
		},
		{
			name: "slices_overlap_total",
			v: func() *[2][]byte {
				backarray := [32]byte{}
				return &[2][]byte{backarray[:], backarray[:]}
			}(),
			want: 2*sizeofSlice + 32,
		},
		{
			name: "slices_overlap",
			v: func() *[4][]uint16 {
				backarray := [32]uint16{}
				return &[4][]uint16{
					backarray[2:4],
					backarray[10:12],
					backarray[20:25],
					backarray[:],
				}
			}(),
			want: 4*sizeofSlice + 32*2,
		},
		{
			name: "slices_overlap_array",
			v: func() *struct {
				a [32]byte
				s [2][]byte
			} {
				v := struct {
					a [32]byte
					s [2][]byte
				}{}
				v.s[0] = v.a[2:4]
				v.s[1] = v.a[5:8]
				return &v
			}(),
			want: 32 + 2*sizeofSlice,
		},
		{
			name: "interface",
			v:    &[2]interface{}{uint64(0), &struct16{}},
			want: 2*sizeofInterface + 8 + 16,
		},
		{
			name: "interface_nil",
			v:    &[2]interface{}{nil, nil},
			want: 2 * sizeofInterface,
		},
		{
			name: "structiface_slice",
			v:    &structiface{x: make([]byte, 10)},
			want: sizeofWord + sizeofInterface + sizeofSlice + 10,
		},
		{
			name: "structiface_pointer",
			v: func() *structiface {
				s := &struct16{1, 2}
				return &structiface{s: s, x: &s.x}
			}(),
			want: sizeofWord + 16 + sizeofInterface,
		},
		{
			name: "empty_chan",
			v: func() *chan uint64 {
				c := make(chan uint64)
				return &c
			}(),
			want: sizeofChan,
		},
		{
			name: "empty_closed_chan",
			v: func() *chan uint64 {
				c := make(chan uint64)
				close(c)
				return &c
			}(),
			want: sizeofChan,
		},
		{
			name: "empty_chan_buffer",
			v: func() *chan uint64 {
				c := make(chan uint64, 10)
				return &c
			}(),
			want: sizeofChan + 10*8,
		},
		{
			name: "chan_buffer",
			v: func() *chan uint64 {
				c := make(chan uint64, 10)
				for i := 0; i < 8; i++ {
					c <- 0
				}
				return &c
			}(),
			want: sizeofChan + 10*8,
		},
		{
			name: "closed_chan_buffer",
			v: func() *chan uint64 {
				c := make(chan uint64, 10)
				for i := 0; i < 8; i++ {
					c <- 0
				}
				close(c)
				return &c
			}(),
			want: sizeofChan + 10*8,
		},
		{
			name: "chan_buffer_escan",
			v: func() *chan *struct16 {
				c := make(chan *struct16, 10)
				for i := 0; i < 8; i++ {
					c <- &struct16{x: uint64(i)}
				}
				return &c
			}(),
			want: sizeofChan + 10*sizeofWord + 8*16,
		},
		{
			name: "closed_chan_buffer_escan",
			v: func() *chan *struct16 {
				c := make(chan *struct16, 10)
				for i := 0; i < 8; i++ {
					c <- &struct16{x: uint64(i)}
				}
				close(c)
				return &c
			}(),
			want: sizeofChan + 10*sizeofWord + 8*16,
		},
		{
			name: "nil_chan",
			v: func() *chan *struct16 {
				var c chan *struct16
				return &c
			}(),
			want: sizeofChan,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			size := Scan(test.v)
			if size.Total != test.want {
				t.Errorf("total=%d, want %d", size.Total, test.want)
				t.Logf("\n%s", size.Report())
			}
		})
	}
}
