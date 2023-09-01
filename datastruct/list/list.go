package list

type consumer func(i int, val any) bool
type List interface {
	Add(val any)
	Get(index int) (val any)
	Set(index int, val any)
	Insert(index int, val any)
	Remove(index int) (val any)
	RemoveLast() (val any)
	RemoveAllByVal(val any) int
	RemoveByVal(val any, count int) int
	ReverseRemoveByVal(val any, count int) int
	Len() int
	ForEach(consumer consumer)
	Contains(val any) int
	Range(start int, stop int) []any
}
