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
			if !reflect.DeepEqual(node.Element.Member, tt.exMember) || !reflect.DeepEqual(node.Element.Score, tt.exScore) {
				t.Errorf("insert() %s err: input: (%v, %v), want: (%v, %v)", tt.name, node.Element.Member, node.Element.Score,
					tt.exMember, tt.exScore)
			}
			rank := skiplist.getRank(tt.member, tt.score)
			if !reflect.DeepEqual(rank, tt.rank) {
				t.Errorf("getRank() %s err: reall: %v, want: %v", tt.name, rank, tt.rank)
			}
		})
	}
}

func Test_skipList_getFirstInScoreRange(t *testing.T) {
	skiplist := makeSkipList()
	testCases := []struct {
		name string
		// input
		member string
		score  float64
		min    ScoreBorder
		max    ScoreBorder

		// expect
		exMember string
		exScore  float64
		rank     int64
	}{
		{
			name:     "test1",
			member:   "a",
			score:    float64(0),
			min:      ScoreBorder{Value: float64(0)},
			max:      ScoreBorder{Value: float64(1)},
			exMember: "a",
			exScore:  float64(0),
			rank:     0,
		}, {
			name:     "test2",
			member:   "b",
			score:    float64(1),
			min:      ScoreBorder{Value: float64(0)},
			max:      ScoreBorder{Value: float64(2)},
			exMember: "a",
			exScore:  float64(0),
			rank:     1,
		}, { // ab会插入到a 和 b之间
			name:     "test3",
			member:   "ab",
			score:    float64(1),
			min:      ScoreBorder{Value: float64(1)},
			max:      ScoreBorder{Value: float64(2)},
			exMember: "ab",
			exScore:  float64(1),
			rank:     1,
		}, {
			name:     "test4",
			member:   "a",
			score:    float64(99),
			min:      ScoreBorder{Value: float64(2)},
			max:      ScoreBorder{Value: float64(100)},
			exMember: "a",
			exScore:  float64(99),
			rank:     3,
		},
	}
	for _, tt := range testCases {
		skiplist.insert(tt.member, tt.score)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			node := skiplist.getFirstInScoreRange(&tt.min, &tt.max)
			if !reflect.DeepEqual(node.Element.Member, tt.exMember) ||
				!reflect.DeepEqual(node.Element.Score, tt.exScore) {
				t.Errorf("%s err: reall: %v, want: %v, %v", tt.name, node, tt.member, tt.score)
			}
		})
	}
}

func Test_skipList_getLastInScoreRange(t *testing.T) {
	skiplist := makeSkipList()
	testCases := []struct {
		name string
		// input
		member string
		score  float64
		min    *ScoreBorder
		max    *ScoreBorder

		// expect
		exMember string
		exScore  float64
		rank     int64
	}{
		{
			name:     "test1",
			member:   "a",
			score:    float64(2),
			min:      &ScoreBorder{Value: float64(0)},
			max:      &ScoreBorder{Value: float64(1)},
			exMember: "b",
			exScore:  float64(1),
			rank:     0,
		}, {
			name:     "test2",
			member:   "b",
			score:    float64(1),
			min:      &ScoreBorder{Value: float64(0)},
			max:      &ScoreBorder{Value: float64(2)},
			exMember: "a",
			exScore:  float64(2),
			rank:     1,
		}, { // ab会插入到a 和 b之间
			name:     "test3",
			member:   "ab",
			score:    float64(1),
			min:      &ScoreBorder{Value: float64(1)},
			max:      &ScoreBorder{Value: float64(2)},
			exMember: "a",
			exScore:  float64(2),
			rank:     1,
		}, {
			name:     "test4",
			member:   "c",
			score:    float64(99),
			min:      &ScoreBorder{Value: float64(2)},
			max:      &ScoreBorder{Value: float64(100)},
			exMember: "c",
			exScore:  float64(99),
			rank:     3,
		},
	}
	for _, tt := range testCases {
		skiplist.insert(tt.member, tt.score)
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			node := skiplist.getLastInScoreRange(tt.min, tt.max)
			if !reflect.DeepEqual(node.Element.Member, tt.exMember) ||
				!reflect.DeepEqual(node.Element.Score, tt.exScore) {
				t.Errorf("%s err: reall: %v, want: %s, %f", tt.name, node, tt.member, tt.score)
			}
		})
	}
}
