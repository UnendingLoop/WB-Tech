package main

import "fmt"

type Shard struct {
	name  string
	start int
	end   int
	limit int
}

func getProductShard(productID int, shards []Shard) Shard {
	result := shardFinder(0, len(shards)-1, productID, shards)
	if result == -1 {
		return Shard{}
	}
	return shards[result]
}

func shardFinder(left, right, productID int, shards []Shard) int {
	middle := (left + right) / 2
	if left > right {
		return -1
	}
	if productID >= shards[middle].start && productID <= shards[middle].end {
		return middle
	}
	if productID < shards[middle].start {
		return shardFinder(left, middle-1, productID, shards)
	}
	if productID > shards[middle].end {
		return shardFinder(middle+1, right, productID, shards)
	}
	return -1
}

func main() {
	/*
		a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		s := a[2:5] //[3,4,5] len=3, cap=8
		fmt.Printf("S len: %d, cap: %d\n", len(s), cap(s))
		s2 := s[5:7] //[8,9] len=2, cap=3
		fmt.Printf("S2 len: %d, cap: %d\n", len(s2), cap(s2))

		s2 = append(s2, 11) //[8,9,11] len=3, cap=3
		fmt.Printf("S2 after 1st append len: %d, cap: %d\n", len(s2), cap(s2))

		s2 = append(s2, 12) //[8,9,11,12] len=4, cap=6 - переаллокация
		fmt.Printf("S2 after 2nd append len: %d, cap: %d\n", len(s2), cap(s2))

		fmt.Println(a) //[1,2,3,4,5,6,7,8,9,11] len=10, cap=10
	*/ ////////////////////////////////////////////////////////////////
	input := []Shard{
		{"1", 0, 20000, 20001},
		{"2", 20001, 40000, 20000},
		{"3", 40001, 60000, 20000},
		{"4", 60001, 80000, 20000},
		{"5", 80001, 90000, 10000},
		{"6", 90001, 100000, 10000},
	}

	targetShard := getProductShard(80555, input)
	fmt.Println(targetShard)
}
