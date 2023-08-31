package list

type expected func(val any) bool
type consumer func(i int, val any) bool
type List interface {
	Add(val any)
	Get(index int) (val any)
	Set(index int, val any)
	Insert(index int, val any)
	Remove(index int) (val any)
	RemoveLast() (val any)
	RemoveAllByVal(expected expected) int
	RemoveByVal(expected expected, count int) int
	ReverseRemoveByVal(expected expected, count int) int
	Len() int
	ForEach(consumer consumer)
	Contains(expected expected) int
	Range(start int, stop int) []any
}
