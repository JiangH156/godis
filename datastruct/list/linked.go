package list

type LinkedList struct {
	first *node
	last  *node
	size  int
}

type node struct {
	data any
	prev *node
	next *node
}

func (list *LinkedList) find(index int) (n *node) {
	if index < list.size/2 {
		n = list.first
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n = list.last
		for i := list.size - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

func (list *LinkedList) Add(val any) {
	if list == nil {
		panic("list is nil")
	}
	n := &node{
		data: val,
	}
	if list.last == nil {
		list.first = n
		list.last = n
	} else {
		n.prev = list.last
		list.last.next = n
		list.last = n
	}
	list.size++
}

func (list *LinkedList) Get(index int) (val any) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of range")
	}
	l := &node{}
	l = list.first
	for i := 0; i < index; i++ {
		l = l.next
	}
	return l.data
}

func (list *LinkedList) Set(index int, val any) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of range")
	}
	n := list.find(index)
	n.data = val
}

func (list *LinkedList) Insert(index int, val any) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of range")
	}
	if index == list.size {
		list.Add(val)
		return
	}
	pivot := list.find(index)
	n := &node{
		data: val,
		prev: pivot.prev,
		next: pivot,
	}
	if pivot.prev == nil {
		list.first = n
	} else {
		pivot.prev.next = n
	}
	pivot.prev = n
	list.size++
}

func (list *LinkedList) removeNode(n *node) {
	if n.prev == nil {
		list.first = n.next
	} else {
		n.next.prev = n.prev
	}
	if n.next == nil {
		list.last = n.prev
	} else {
		n.prev.next = n.next
	}
	n.prev = nil
	n.next = nil
	list.size--
}

func (list *LinkedList) Remove(index int) (val any) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of range")
	}
	n := list.find(index)
	list.removeNode(n)
	return n.data
}

func (list *LinkedList) RemoveLast() (val any) {
	n := list.last
	list.removeNode(n)
	return n.data
}

func (list *LinkedList) RemoveAllByVal(expected expected) int {
	if list == nil {
		panic("list is nil")
	}
	n := list.first
	removed := 0
	var nextNode *node
	for n != nil {
		nextNode = n.next
		if expected(n) {
			list.removeNode(n)
			removed++
		}
		n = nextNode
	}
	return removed
}

func (list *LinkedList) RemoveByVal(expected expected, count int) int {
	if list == nil {
		panic("list is nil")
	}
	n := list.first
	removed := 0
	var nextNode *node
	for n != nil && removed < count {
		nextNode = n.next
		if expected(n) {
			list.removeNode(n)
			removed++
		}
		n = nextNode
	}
	return removed
}

func (list *LinkedList) ReverseRemoveByVal(expected expected, count int) int {
	if list == nil {
		panic("list is nil")
	}
	n := list.last
	removed := 0
	var prevNode *node
	for n != nil && removed < count {
		prevNode = n.prev
		if expected(n) {
			list.removeNode(n)
			removed++
		}
		n = prevNode
	}
	return removed
}

func (list *LinkedList) Len() int {
	if list == nil {
		panic("list is nil")
	}
	return list.size
}

func (list *LinkedList) ForEach(consumer consumer) {
	if list == nil {
		panic("list is nil")
	}
	i := 0
	n := list.first
	var nextNode *node
	for n != nil {
		nextNode = n.next
		ok := consumer(i, n.data)
		if !ok {
			break
		}
		i++
		n = nextNode
	}
}

func (list *LinkedList) Contains(expected expected) int {
	if list == nil {
		panic("list is nil")
	}
	result := 0
	n := list.first
	var nextNode *node
	for n != nil {
		nextNode = n.next
		if expected(n) {
			result++
		}
		n = nextNode
	}
	return result
}

func (list *LinkedList) Range(start int, stop int) (vals []any) {
	if list == nil {
		panic("list is nil")
	}
	if start > stop {
		panic("start cannot be greater than stop")
	}
	if start < 0 {
		start = 0
	}
	if stop > list.size-1 {
		stop = list.size - 1
	}

	i := 0
	n := list.first
	for n != nil {
		if i >= start {
			vals = append(vals, n.data)
		}
		if i > stop {
			break
		}
		i++
		n = n.next
	}
	return vals
}
func Make(vals ...any) *LinkedList {
	list := &LinkedList{}
	for _, val := range vals {
		list.Add(val)
	}
	return list
}
