package set

import "github.com/jiangh156/godis/datastruct/dict"

type Set struct {
	dict dict.Dict
}

func Make() *Set {
	return &Set{
		dict: dict.MakeSyncDict(),
	}
}

func MakeFromVals(members ...string) *Set {
	set := &Set{
		dict: dict.MakeSyncDict(),
	}
	return set
}

func (set *Set) Add(member string) int {
	if set == nil {
		panic("set is nil")
	}
	return set.dict.Put(member, struct{}{})
}
func (set *Set) Remove(member string) int {
	if set == nil {
		panic("set is nil")
	}
	return set.dict.Remove(member)
}
func (set *Set) Has(member string) bool {
	if set == nil {
		panic("set is nil")
	}
	_, exists := set.dict.Get(member)
	return exists
}
func (set *Set) Len() int {
	if set == nil {
		panic("set is nil")
	}
	return set.dict.Len()
}
func (set *Set) ToSlice() []string {
	if set == nil {
		panic("set is nil")
	}
	slice := make([]string, set.Len())
	i := 0
	set.dict.ForEach(func(member string, value any) bool {
		if i < len(slice) {
			slice[i] = member
		} else {
			slice = append(slice, member)
		}
		return true
	})
	return slice
}
func (set *Set) ForEach(consumer func(member string) bool) {
	set.dict.ForEach(func(key string, value any) bool {
		return consumer(key)
	})
}
func (set *Set) Union(another *Set) *Set {
	if set == nil {
		panic("set is nil")
	}
	result := Make()
	another.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	set.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	return result
}
func (set *Set) Diff(another *Set) *Set {
	if set == nil {
		panic("set is nil")
	}
	result := Make()
	another.ForEach(func(member string) bool {
		if !another.Has(member) {
			result.Add(member)
		}
		return true
	})
	return result
}
func (set *Set) RandomMembers(limit int) []string {
	return set.dict.RandomKeys(limit)
}
func (set *Set) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
