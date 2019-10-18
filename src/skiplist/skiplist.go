/**
 *Create by chensr on 2019/10/18
*/

package skiplist

import "math/rand"

type ObjInterface interface {
	GetKey() string
	IsValid() bool
	Compare(obj ObjInterface) int
}

type ScoreInterface interface {
	Compare(obj ScoreInterface) int
	GetValue() ScoreInterface
	IsValid() bool
}

type Node struct {
	Obj ObjInterface
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
	RandNum int
}

func CreateSkipList(maxLevel int, randNum int) *SkipList {
	sl := new(SkipList)

	sl.LevelNum = 1
	sl.Length = 0
	sl.Header = sl.CreateNode(0, nil, nil)
	sl.RandNum = randNum
	for i := 0; i < maxLevel; i++ {
		sl.Header.LevelInfo[i].Forward = nil
		sl.Header.LevelInfo[i].Span = 0
	}

	return sl
}

func (sl *SkipList) CreateNode(level int, obj ObjInterface, score ScoreInterface) *Node {
	node := new(Node)
	node.LevelInfo = make([]*Level, level)
	node.Obj = obj
	node.Score = score
	return node
}

func (sl *SkipList) Insert(obj ObjInterface, score ScoreInterface) {
	if !obj.IsValid() || !score.IsValid() {
		return
	}
	var rank [sl.LevelNum]int
	var update [sl.LevelNum]*Node
	x := sl.Header

	for i := sl.LevelNum-1; i >= 0; i-- {
		if i == sl.LevelNum - 1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for {
			compareScore := score.Compare(x.LevelInfo[i].Forward.Score)
			compareKey := obj.Compare(x.LevelInfo[i].Forward.Obj)
			if x.LevelInfo[i].Forward != nil && (
				compareScore == 1 || (compareScore == 0 && compareKey == 1)) {
				rank[i] += x.LevelInfo[i].Span
				x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}
		update[i] = x
	}

	level := rand.Intn(sl.RandNum)
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
		x.LevelInfo[i] .Forward = update[i].LevelInfo[i].Forward
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

func (sl *SkipList)GetRank(obj ObjInterface, score ScoreInterface) int {
	var x *Node
	var rank int

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i-- {
		for {
			scoreCompare := score.Compare(x.LevelInfo[i].Forward.Score)
			objCompare := obj.Compare(x.LevelInfo[i].Forward.Obj)
			if x.LevelInfo[i].Forward != nil && (
				scoreCompare == 1 || (scoreCompare == 0 && objCompare >= 0)) {
				rank += x.LevelInfo[i].Span
				x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}

		if x.Obj != nil && obj.Compare(x.Obj) == 0 {
			return rank
		}
	}

	return 0
}

func (sl *SkipList) GetObjByRank(rank int) *Node {
	var x *Node
	var traversed int

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i++ {
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

func (sl *SkipList) Delete(obj ObjInterface, score ScoreInterface) bool {
	update := make([]*Node, sl.LevelNum)
	var x *Node

	x = sl.Header
	for i := sl.LevelNum - 1; i >= 0; i-- {
		for {
			scoreCompare := score.Compare(x.LevelInfo[i].Forward.Score)
			objCompare := obj.Compare(x.LevelInfo[i].Forward.Obj)
			if x.LevelInfo[i].Forward != nil && (
				scoreCompare == 1 || (scoreCompare == 0 && objCompare > 0)) {
					x = x.LevelInfo[i].Forward
			} else {
				break
			}
		}
		update[i] = x
	}

	x = x.LevelInfo[0].Forward
	if x != nil {
		scoreCompare := score.Compare(x.Score)
		objCompare := obj.Compare(x.Obj)
		if scoreCompare == 0 && objCompare == 0 {
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