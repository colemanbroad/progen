package main

import (
	"math"
)

var reward_buckets map[int]float64
var reward []RewardTime
var Reward_total float64
var global_time int

type RewardTime struct {
	Reward float64
	Time   int
	Source string
}

func init_reward() {
	reward_buckets = make(map[int]float64)
	reward = make([]RewardTime, 0)
	Reward_total = 0
	global_time = 0
}

func Found_reward(r float64, source string) {
	Reward_total += r
	reward = append(reward, RewardTime{Reward: r, Time: global_time, Source: source})
	// fmt.Println("found reward: ", r)
}

func addPowerOfTwo() {
	addFuncToLibrary(isPowerOfTwo, "isPowerOfTwo", []Type{"int"}, "bool")
}

func isPowerOfTwo(n int) bool {
	c, exists := reward_buckets[n]
	if !exists {
		c = 0.0
	}
	reward_buckets[n] = c + 1.0

	d := math.Log2(float64(n))
	r := 0.0
	is_power := false
	if math.Floor(d) == d {
		r = 1.0 / (c + 1)
		// r = 1.0
		Found_reward(r, "isPowerOfTwo")
		is_power = true
	}
	history_power_of_two = append(history_power_of_two, Reward_power_of_two{Value: float32(n), Reward: float32(r), Time: global_time})
	// fmt.Println("pow2 ", n, d, is_power)

	return is_power
}

func addPeanoLib() {
	var f any
	f = func(a int) int { return a + 1 }
	addFuncToLibrary(f, "succ", []Type{"int"}, "int")
	f = func() int { return 0 }
	addFuncToLibrary(f, "zero", []Type{}, "int")
}
