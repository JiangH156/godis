package sortedset

import "math/rand"

const (
	maxLevel int16 = 16
)

type Element struct {
	Member string
	Score  float64
}

type Node struct {
	Element  Element  //元素的名称和score
	backward *Node    // 后向指针
	level    []*Level // 垂直层，Level[0] 为最下层
}

type Level struct {
	forward *Node // 指向同层中的下一个节点
	span    int64 // 到forward跳过的节点数
}

type skipList struct {
	header *Node
	tail   *Node
	length int64
	level  int16
}

func makeNode(level int16, score float64, member string) *Node {
	n := &Node{
		Element: Element{
			Score:  score,
			Member: member,
		},
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = new(Level)
	}
	return n
}

func makeSkipList() *skipList {
	return &skipList{
		level:  1,
		header: makeNode(maxLevel, 0, ""),
		//header: nil,
	}
}
func randomLevel() int16 {
	level := int16(1)
	for float32(rand.Int31()&0xFFFF) < 0.25*0xFFFF {
		level++
	}
	if level < maxLevel {
		return level
	}
	return maxLevel
}
func (skipList *skipList) insert(member string, score float64) *Node {
	update := make([]*Node, maxLevel)
	rank := make([]int64, maxLevel) // rank[i] 记录 第i个节点与header的距离

	// 寻找先驱节点
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历
		// 当i为第一个遍历到的
		if i == skipList.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		// 同一level下不断寻找节点
		if n.level[i] != nil {
			// 遍历搜索
			for n.level[i].forward != nil && (n.level[i].forward.Element.Score < score ||
				(n.level[i].forward.Element.Score == score && n.level[i].forward.Element.Member < member)) { // same score, different member
				rank[i] += n.level[i].span
				n = n.level[i].forward
			}
		}
		update[i] = n
	}

	// 节点层数
	level := randomLevel()
	// 当新插入的节点的层数大于当前跳表的最大长度时，更新最大长度
	if level > skipList.level {
		for i := skipList.level; i < level; i++ {
			rank[i] = 0
			update[i] = skipList.header
			update[i].level[i].span = skipList.length // 这里并不是等于length，后面会更新
		}
		skipList.level = level
	}

	// 插入数据
	// 新建节点
	node := makeNode(level, score, member)
	for i := int16(0); i < level; i++ {
		// 新节点的 forward 指向先驱节点的 forward
		node.level[i].forward = update[i].level[i].forward
		// 先驱节点的 forward 指向新节点
		update[i].level[i].forward = node
		// 计算先驱节点和新节点的 span
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = rank[0] - rank[i] + 1
	}

	// 考虑特殊情况：表头、表尾
	// 考虑当前节点的backward,只需考虑前
	if update[0] == nil {
		node.backward = skipList.header
	} else {
		node.backward = update[0]
	}

	//考虑当前节点的forward->backward
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node
	} else {
		skipList.tail = node
	}
	skipList.length++
	return node
}

func (skipList *skipList) getRank(member string, score float64) int64 {
	var rank int64 = 0
	// 寻找先驱节点
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历
		// 同一level下不断寻找节点
		if n.level[i] != nil {
			// 找到元素位置
			if n.Element.Score == score && n.Element.Member == member {
				break
			}
			// 同一层次遍历
			for n.level[i].forward != nil && (n.level[i].forward.Element.Score < score ||
				(n.level[i].forward.Element.Score == score && n.level[i].forward.Element.Member < member)) { // same score, different member
				rank += n.level[i].span
				n = n.level[i].forward
			}
		}
	}
	return rank
}

// 寻找排名为 rank 的节点
func (skipList *skipList) getByRank(rank int64) *Node {
	var r int64 = -1
	// 寻找先驱节点
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历

		if n.level[i] != nil {
			// 同一层次遍历
			for n.level[i].forward != nil && (r+n.level[i].span <= rank) { //同一层中存在后续节点
				// 同一level下不断寻找节点
				if r == rank {
					return n
				}
				r += n.level[i].span
				n = n.level[i].forward
			}
		}
	}
	return n
}
func (skipList *skipList) hasInRange(min *ScoreBorder, max *ScoreBorder) bool {
	if min.Value > max.Value || (min.Value == max.Value && (min.Exclude || max.Exclude)) {
		return false
	}
	if min.Value > skipList.tail.Element.Score {
		return false
	}
	if max.Value < skipList.header.level[0].forward.Element.Score {
		return false
	}
	return true
}

// finish test
func (skipList *skipList) getFirstInScoreRange(min *ScoreBorder, max *ScoreBorder) *Node {
	if !skipList.hasInRange(min, max) {
		return nil
	}
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历
		// 同一level下不断寻找节点
		if n.level[i] != nil {
			// 同一层次遍历
			for n.level[i].forward != nil && max.greater(n.level[i].forward.Element.Score) {
				n = n.level[i].forward
				// 找到元素位置
				if min.less(n.Element.Score) { // 大于最小值
					return n
				}
			}
		}
	}
	return n
}

func (skipList *skipList) getLastInScoreRange(min *ScoreBorder, max *ScoreBorder) *Node {
	if !skipList.hasInRange(min, max) {
		return nil
	}
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历
		// 同一level下不断寻找节点
		if n.level[i] != nil {
			// 同一层次遍历
			for n.level[i].forward != nil && max.greater(n.level[i].forward.Element.Score) {
				n = n.level[i].forward
			}
		}
	}
	return n
}
func (skipList *skipList) removeRangeByScore(min *ScoreBorder, max *ScoreBorder) (removed []*Element) {
	if !skipList.hasInRange(min, max) {
		return nil
	}
	reNodes := []*Element{}

	var n = skipList.header
	for i := skipList.length - 1; i >= 0; i-- {
		if n.level[i] != nil {
			for n.level[i].forward != nil && min.Value > n.level[i].forward.Element.Score {
				n = n.level[i].forward
			}
		}
	}

	for n != nil {
		if n.level[0] != nil {
			if !max.greater(n.Element.Score) {
				break
			}
			member := n.Element.Member
			score := n.Element.Score
			n = n.level[0].forward
			if min.less(score) && max.greater(score) {
				reNodes = append(reNodes, &Element{Member: member, Score: score})
				skipList.remove(member, score)
			}
		}
	}
	return reNodes
}
func (skipList *skipList) removeRangeByRank(start int64, stop int64) (removed []*Element) {
	reNodes := []*Element{}
	if start > skipList.length {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if stop > skipList.length {
		stop = skipList.length
	}
	var n = skipList.header
	var rank int64 = 0
	for i := skipList.length - 1; i >= 0; i-- {
		if n.level[i] != nil {
			for n.level[i].forward != nil && rank+n.level[i].span < start {
				rank += n.level[i].span
				n = n.level[i].forward
			}
		}
	}
	i := int64(0)
	l := stop - start - 1
	for i < l {
		if n.level[0].forward != nil {
			n = n.level[0].forward
			reNodes = append(reNodes, &Element{Member: n.Element.Member, Score: n.Element.Score})
			skipList.remove(n.Element.Member, n.Element.Score)
			i++
		}
	}
	return reNodes
}

// finish test
func (skipList *skipList) removeNode(node *Node, update []*Node) {
	var level = len(update)
	// 更新span
	for i := level - 1; i >= 0; i-- {
		if update[i] != nil {
			if update[i].level[i].forward == node {
				update[i].level[i].span += node.level[i].span - 1
				update[i].level[i].forward = node.level[i].forward
			} else {
				update[i].level[i].span--
			}
		}
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = update[0]
	} else {
		skipList.tail = update[0]
	}
	skipList.length--
}

// finish test
func (skipList *skipList) remove(member string, score float64) bool {
	update := make([]*Node, maxLevel)
	// 寻找先驱节点
	// 当前遍历节点
	var n = skipList.header
	// 从上往下遍历
	for i := skipList.level - 1; i >= 0; i-- { // 自顶向下遍历

		// 同一level下不断寻找节点
		if n.level[i] != nil {
			// 遍历搜索
			for n.level[i].forward != nil && (n.level[i].forward.Element.Score < score ||
				(n.level[i].forward.Element.Score == score && n.level[i].forward.Element.Member < member)) { // same score, different member
				n = n.level[i].forward
			}
		}
		update[i] = n
	}
	if n.level[0].forward == nil || n.level[0].forward.Element.Member != member || n.level[0].forward.Element.Score != score {
		return false
	}
	skipList.removeNode(n.level[0].forward, update)
	return true
}
