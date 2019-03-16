package main

import (
	"flag"
	"fmt"
	"math/rand"
)

// maxGames is a safety limit to ensure we don't run in an infinite loop.
const maxGames = 10000

var rankMilestone = map[int]struct{}{
	20: {},
	15: {},
	10: {},
	5: {},
}

func main() {
	var (
		startingRank = flag.Int("rank", 25, "At which rank to start the simulation.")
		startingStar = flag.Int("star", 0, "At which star to start the simulation.")
		winRate = flag.Float64("rate", 0.5, "Estimated deck win rate.")
		runs = flag.Int("runs", 1000, "Number of runs.")
	)

	flag.Parse()

	results := make([]int, 0, *runs)

	for i := 0; i < *runs; i++ {
		rank := *startingRank
		stars := *startingStar
		winStreak := 0
		games := 0

		for rank > 0 && games < maxGames {
			games++

			if simulateMatch(*winRate) {
				winStreak++
				rank, stars = winStar(rank, stars, winStreak)
			} else {
				winStreak = 0
				rank, stars = looseStar(rank, stars)
			}
		}

		results = append(results, games)
	}

	min, max, avg := stats(results)

	fmt.Printf("estimated number of matches: %.2f (min: %d, max: %d)\n", avg, min, max)
}

func simulateMatch(winRate float64) bool {
	return rand.Float64() < winRate
}

func winStar(rank int, stars int, winStreak int) (int, int) {
	stars += earnStar(rank, winStreak)
	if stars > rankStars(rank) {
		stars -= rankStars(rank)
		rank--
	}
	return rank, stars
}

func looseStar(rank, stars int) (int, int) {
	stars--
	if stars < 0 {
		if !milestone(rank) {
			rank++
		}
		stars = rankStars(rank)
	}
	return rank, stars
}

func earnStar(rank int, winStreak int) int {
	if rank > 5 && winStreak >= 3 {
		return 2
	}
	return 1
}

func rankStars(rank int) int {
	if rank >= 50 && rank < 15 {
		return 3
	}
	if rank >= 15 && rank < 10 {
		return 4
	}
	return 5
}

func milestone(rank int) bool {
	_, ok := rankMilestone[rank]
	return ok
}

func stats(results []int) (int, int, float64) {
	if len(results) == 0 {
		return 0, 0, .0
	}

	total := results[0]
	min := results[0]
	max := results[0]

	for i := 1; i < len(results); i++ {
		total += results[i]
		if results[i] < min {
			min = results[i]
		}
		if results[i] > max {
			max = results[i]
		}
	}

	return min, max, float64(total) / float64(len(results))
}
