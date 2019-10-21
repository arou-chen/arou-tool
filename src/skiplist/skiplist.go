/**
 *Create by chensr on 2019/10/18
*/

package skiplist

import "math/rand"

//type ObjInterface interface {
//	GetKey() string
//	IsValid() bool
//	Compare(obj ObjInterface) int
//}

type ScoreInterface interface {
	Compare(obj interface{}) int
	//GetValue() ScoreInterface
	//IsValid() bool
}

type Node struct {
	Obj string
	Score ScoreInterface
	BackWard *Node
	LevelInfo []*Level
}

type Level struct {
	Forward *Node
	Span int
}

type SkipList struct {
	Header *Node
	Tail *Node
	Length int
	LevelNum int
	MaxLevel int
}

func CreateSkipList(maxLevel int) *SkipList {
	sl := new(SkipList)

	sl.LevelNum = 1
	sl.Length = 0
	sl.Header = sl.CreateNode(0, "", nil)
	sl.MaxLevel = maxLevel
	sl.Header.LevelInfo = make([]*Level, maxLevel)
	for i := 0; i < maxLevel; i++ {
		sl.Header.LevelInfo[i] = &Level{}
	}

	return sl
}

func (sl *SkipList) CreateNode(level int, obj string, score ScoreInterface) *Node {
	node := new(Node)
	node.LevelInfo = make([]*Level, level)
	node.Obj = obj
	node.Score = score
	return node
}

func (sl *SkipList) Insert(obj string, score ScoreInterface) {
	if obj == "" {
		return
	}
	rank := make([]int, sl.MaxLevel)
	update := make([]*Node, sl.MaxLevel)
	x := sl.Header

	for i := sl.LevelNum-1; i >= 0; i-- {
		if i == sl.LevelNum - 1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for {
			//compareScore := score.Compare(x.LevelInfo[i].Forward.Score)
			//compareKey := obj >= x.LevelInfo[i].Forward.Obj
			if x.LevelInfo[i].Forward != nil && (
				score.Compare(x.LevelInfo[i].Forward.Score) == 1 ||
					(score.Compare(x.LevelInfo[i].Forward.Score) == 0 && obj > x.LevelInfo[i].Forward.Obj)) {
				rank[i] += x.LevelInfo[i].Span
				x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}
		update[i] = &Node{}
		update[i] = x
	}

	level := rand.Intn(sl.MaxLevel)
	if level > sl.LevelNum {
		for i := sl.LevelNum; i < level; i++ {
			rank[i] = 0
			update[i] = sl.Header
			update[i].LevelInfo[i].Span = sl.Length
		}
		sl.LevelNum = level
	}

	x = sl.CreateNode(level, obj, score)
	for i := 0; i < level; i++ {
		x.LevelInfo[i] = &Level{}
		x.LevelInfo[i].Forward = update[i].LevelInfo[i].Forward
		update[i].LevelInfo[i].Forward = x
		randDelta := rank[0]-rank[i]
		x.LevelInfo[i].Span = update[i].LevelInfo[i].Span - (randDelta)
		update[i].LevelInfo[i].Span = randDelta + 1
	}

	for i := level; i < sl.LevelNum; i++ {
		update[i].LevelInfo[i].Span++
	}

	if update[0] == sl.Header {
		x.BackWard = nil
	} else {
		x.BackWard = update[0]
	}

	if x.LevelInfo[0].Forward != nil {
		x.LevelInfo[0].Forward.BackWard = x
	} else {
		sl.Tail = x
	}

	sl.Length++
}

func (sl *SkipList)GetRank(obj string, score ScoreInterface) int {
	var x *Node
	var rank int

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i-- {
		for {
			if x.LevelInfo[i].Forward != nil && (
				score.Compare(x.LevelInfo[i].Forward.Score) == 1 ||
					(score.Compare(x.LevelInfo[i].Forward.Score) == 0 && obj >= x.LevelInfo[i].Forward.Obj)) {
				rank += x.LevelInfo[i].Span
				x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}

		if x.Obj != "" && obj == x.Obj {
			return rank
		}
	}

	return 0
}

func (sl *SkipList) GetObjByRank(rank int) *Node {
	var x *Node
	var traversed int

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i-- {
		for {
			if x.LevelInfo[i].Forward != nil && (traversed + x.LevelInfo[i].Span) <= rank {
				traversed += x.LevelInfo[i].Span
				x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}

		if traversed == rank {
			return x
		}
	}

	return nil
}

func (sl *SkipList) Delete(obj string, score ScoreInterface) bool {
	update := make([]*Node, sl.LevelNum)
	var x *Node

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i-- {
		for {
			if x.LevelInfo[i].Forward != nil && (
				score.Compare(x.LevelInfo[i].Forward.Score) == 1 ||
					(score.Compare(x.LevelInfo[i].Forward.Score) == 0 && obj > x.LevelInfo[i].Forward.Obj)) {
					x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}
		update[i] = &Node{}
		update[i] = x
	}

	x = x.LevelInfo[0].Forward
	if x != nil {
		scoreCompare := score.Compare(x.Score)
		if scoreCompare == 0 && obj == x.Obj {
			sl.DeleteNode(x, update)
			//arou 会不会有内存释放问题
			x = &Node{}
			return true
		}
	}

	return false
}

func (sl *SkipList) DeleteNode(x *Node, update []*Node) {
	for i := 0; i < sl.LevelNum; i++ {
		if update[i].LevelInfo[i].Forward == x {
			update[i].LevelInfo[i].Span += x.LevelInfo[i].Span - 1
			update[i].LevelInfo[i].Forward = x.LevelInfo[i].Forward
		} else {
			update[i].LevelInfo[i].Span -= 1
		}
	}

	if x.LevelInfo[0].Forward != nil {
		x.LevelInfo[0].Forward.BackWard = x.BackWard
	} else {
		sl.Tail = x.BackWard
	}

	for {
		if sl.LevelNum > 1 && sl.Header.LevelInfo[sl.LevelNum - 1].Forward == nil {
			sl.LevelNum--
		} else {
			break
		}
	}
	sl.Length--
}