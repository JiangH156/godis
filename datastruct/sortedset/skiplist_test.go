package sortedset

import (
	"reflect"
	"testing"
)

func Test_skipList_insert(t *testing.T) {
	skiplist := makeSkipList()
	testCases := []struct {
		name string
		// input
		member string
		score  float64

		// expect
		exMember string
		exScore  float64
		rank     int64
	}{
		{
			name:     "test1",
			member:   "a",
			score:    float64(1),
			exMember: "a",
			exScore:  float64(1),
			rank:     0,
		}, {
			name:     "test2",
			member:   "b",
			score:    float64(1),
			exMember: "b",
			exScore:  float64(1),
			rank:     1,
		}, { // ab会插入到a 和 b之间
			name:     "test3",
			member:   "ab",
			score:    float64(1),
			exMember: "ab",
			exScore:  float64(1),
			rank:     1,
		}, {
			name:     "test4",
			member:   "a",
			score:    float64(99),
			exMember: "a",
			exScore:  float64(99),
			rank:     3,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			node := skiplist.insert(tt.member, tt.score)
			rank := skiplist.getRank(tt.member, tt.score)
			if !reflect.DeepEqual(node.Element.Member, tt.exMember) || !reflect.DeepEqual(node.Element.Score, tt.exScore) {
				t.Errorf("insert() %s err: input: (%v, %v), want: (%v, %v)", tt.name, node.Element.Member, node.Element.Score,
					tt.exMember, tt.exScore)
			}
			if !reflect.DeepEqual(rank, tt.rank) {
				t.Errorf("getRank() %s err: reall: %v, want: %v", tt.name, rank, tt.rank)
			}
		})
	}
}
