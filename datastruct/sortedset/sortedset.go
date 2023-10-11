package sortedset

import "strconv"

type SortedSet struct {
	dict     map[string]*Element
	skiplist *skipList
}

func Make() *SortedSet {
	return &SortedSet{
		dict:     make(map[string]*Element),
		skiplist: makeSkipList(),
	}
}

/*
 * return: has inserted new node
 */
func (sortedSet *SortedSet) Add(member string, score float64) bool {
	// update dict
	element, ok := sortedSet.dict[member]
	sortedSet.dict[member] = &Element{
		Member: member,
		Score:  score,
	}
	if ok {
		if element.Score != score {
			sortedSet.skiplist.remove(member, element.Score)
			sortedSet.skiplist.insert(member, score)
		} else {
			return false
		}
	} else {
		sortedSet.skiplist.remove(member, element.Score)
		sortedSet.skiplist.insert(member, score)
	}
	return true
}

func (sortedSet *SortedSet) Len() int64 {
	return int64(len(sortedSet.dict))
}

func (sortedSet *SortedSet) Get(member string) (element *Element, ok bool) {
	element, ok = sortedSet.dict[member]
	return
}

func (sortedSet *SortedSet) Remove(member string) bool {
	element, ok := sortedSet.dict[member]
	if !ok {
		return false
	}
	remove := sortedSet.skiplist.remove(member, element.Score)
	if !remove {
		return false
	}
	delete(sortedSet.dict, member)
	return true
}

/**
 * get 0-based rank
 */
func (sortedSet *SortedSet) GetRank(member string, desc bool) (rank int64) {
	element, ok := sortedSet.dict[member]
	if !ok {
		return -1
	}
	r := sortedSet.skiplist.getRank(member, element.Score)
	if desc {
		r = sortedSet.skiplist.length - r
	} else {
		r--
	}
	return r
}

/**
 * traverse [start, stop), 0-based rank
 */
// TODO start和stop边界问题，都是从0开始
func (sortedSet *SortedSet) ForEach(start int64, stop int64, desc bool, consumer func(element *Element) bool) {
	size := sortedSet.Len()
	if start < 0 || start >= size {
		panic("illegal start " + strconv.FormatInt(start, 10))
	}
	if stop < start || stop > size {
		panic("illegal end " + strconv.FormatInt(stop, 10))
	}
	var node *Node
	if desc {
		node = sortedSet.skiplist.tail
		if start > 0 {
			node = sortedSet.skiplist.getByRank(int64(size - start - 1)) //TODO
		}
	} else {
		node = sortedSet.skiplist.header.level[0].forward
		if start > 0 {
			node = sortedSet.skiplist.getByRank(int64(start - 1)) // TODO
		}
	}

	sliceSize := stop - start - 1
	for i := int64(0); i < sliceSize; i++ {
		if !consumer(&node.Element) {
			break
		}
		if desc {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
	}
}

/**
 * return [start, stop), 0-based rank
 * assert start in [0, size), stop in [start, size]
 */
func (sortedSet *SortedSet) Range(start int64, stop int64, desc bool) []*Element {
	size := sortedSet.Len()
	if start < 0 || start >= size {
		return nil
	}
	if stop < start || stop > size {
		return nil
	}
	sliceSize := int(stop - start)
	slice := make([]*Element, sliceSize)
	i := 0
	sortedSet.ForEach(start, stop, desc, func(element *Element) bool {
		slice[i] = element
		i++
		return true
	})
	return slice
}

func (sortedSet *SortedSet) Count(min *ScoreBorder, max *ScoreBorder) int64 {
	var i = int64(0)
	sortedSet.ForEach(0, sortedSet.Len(), false, func(element *Element) bool {
		if !min.less(element.Score) {
			return true
		}
		if !max.greater(element.Score) {
			return false
		}
		i++
		return true
	})
	return i
}

func (sortedSet *SortedSet) ForEachByScore(min *ScoreBorder, max *ScoreBorder, offset int64, limit int64, desc bool, consumer func(element *Element) bool) {
	var node *Node
	if desc {
		node = sortedSet.skiplist.getLastInScoreRange(min, max)
	} else {
		node = sortedSet.skiplist.getFirstInScoreRange(min, max)
	}
	for node != nil && offset != 0 {
		if desc {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
		offset--
	}
	// A negative limit returns all elements from the offset
	for i := 0; (i < int(limit) || limit < 0) && node != nil; i++ {
		if !consumer(&node.Element) {
			break
		}
		if desc {
			node = node.backward
		} else {
			node = node.level[0].forward
		}
		if !min.less(node.Element.Score) || !max.greater(node.Element.Score) {
			break
		}
	}
}

/*
 * param limit: <0 means no limit
 */
func (sortedSet *SortedSet) RangeByScore(min *ScoreBorder, max *ScoreBorder, offset int64, limit int64, desc bool) []*Element {
	if limit == 0 || offset < 0 {
		return make([]*Element, 0)
	}
	elements := make([]*Element, 0)
	sortedSet.ForEachByScore(min, max, offset, limit, desc, func(element *Element) bool {
		elements = append(elements, element)
		return true
	})
	return elements
}

func (sortedSet *SortedSet) RemoveByScore(min *ScoreBorder, max *ScoreBorder) int64 {
	removed := sortedSet.skiplist.removeRangeByScore(min, max)
	for _, element := range removed {
		delete(sortedSet.dict, element.Member)
	}
	return int64(len(removed))
}

/*
 * 0-based rank, [start, stop)
 */
func (sortedSet *SortedSet) RemoveByRank(start int64, stop int64) int64 {
	removed := sortedSet.skiplist.removeRangeByRank(start, stop)
	for _, element := range removed {
		delete(sortedSet.dict, element.Member)
	}
	return int64(len(removed))
}
