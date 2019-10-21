/**
 *Create by chensr on 2019/10/18
 */

package main

import (
	"fmt"
	"skiplist"
	"strconv"
)

type RankData struct {
	UserID     int64
	Epoch      int32
	TechNum    int32
	Population int64
}

func (data *RankData) Compare(inter interface{}) int {
	obj, _ := inter.(*RankData)
	if obj.Epoch != data.Epoch {
		if data.Epoch > obj.Epoch {
			return 1
		} else {
			return -1
		}
	} else if data.TechNum != obj.TechNum {
		if data.TechNum > obj.TechNum {
			return 1
		} else {
			return -1
		}
	} else if data.Population != obj.Population {
		if data.Population > obj.Population {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

var rank *skiplist.SkipList

func main() {
	rank = skiplist.CreateSkipList(10)
	temp1 := &RankData{1, 2, 2, 1000}
	temp2 := &RankData{2, 2, 1, 1000}
	rank.Insert(strconv.Itoa(int(temp1.UserID)), temp1)
	rank.Insert(strconv.Itoa(int(temp2.UserID)), temp2)
	fmt.Println(rank.GetRank(strconv.Itoa(int(temp1.UserID)), temp1))
	fmt.Println(rank.GetObjByRank(1).Score)
	rank.Delete(strconv.Itoa(int(temp2.UserID)), temp2)
	fmt.Println(rank.GetObjByRank(1).Score)
}
