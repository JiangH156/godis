package sortedset

type SortedSet struct {
	dict     map[string]*Element
	skipList *skipList
}
