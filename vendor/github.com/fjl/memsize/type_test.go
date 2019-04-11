package memsize

import (
	"reflect"
	"testing"
)

var typCacheTests = []struct {
	val  interface{}
	want typInfo
}{
	{
		val:  int(0),
		want: typInfo{isPointer: false, needScan: false},
	},
	{
		val:  make(chan struct{}, 1),
		want: typInfo{isPointer: true, needScan: true},
	},
	{
		val:  struct{ A int }{},
		want: typInfo{isPointer: false, needScan: false},
	},
	{
		val:  struct{ S string }{},
		want: typInfo{isPointer: false, needScan: true},
	},
	{
		val:  structloop{},
		want: typInfo{isPointer: false, needScan: true},
	},
	{
		val:  [3]int{},
		want: typInfo{isPointer: false, needScan: false},
	},
	{
		val:  [3]struct{ A int }{},
		want: typInfo{isPointer: false, needScan: false},
	},
	{
		val:  [3]struct{ S string }{},
		want: typInfo{isPointer: false, needScan: true},
	},
	{
		val:  [3]structloop{},
		want: typInfo{isPointer: false, needScan: true},
	},
	{
		val: struct {
			a [32]uint8
			s [2][]uint8
		}{},
		want: typInfo{isPointer: false, needScan: true},
	},
}

func TestTypeInfo(t *testing.T) {
	// This cache is shared among all test cases. It is used
	// to verify that putting many different types into the cache
	// doesn't change the resulting info.
	sharedtc := make(typCache)

	for i := range typCacheTests {
		test := typCacheTests[i]
		typ := reflect.TypeOf(test.val)
		t.Run(typ.String(), func(t *testing.T) {
			tc := make(typCache)
			info := tc.info(typ)
			if !reflect.DeepEqual(info, test.want) {
				t.Fatalf("wrong info from local cache:\ngot %+v, want %+v", info, test.want)
			}
			info = sharedtc.info(typ)
			if !reflect.DeepEqual(info, test.want) {
				t.Fatalf("wrong info from shared cache:\ngot %+v, want %+v", info, test.want)
			}
		})
	}
}
