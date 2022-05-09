package core

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/jordi-reinsma/bagel/db"
	"github.com/jordi-reinsma/bagel/model"
)

func GenerateMatches(dB db.DB, users []model.User) ([]model.Match, error) {
	// sort the users by id
	byID := func(i, j int) bool {
		return users[i].ID < users[j].ID
	}
	sort.Slice(users, byID)
	// make the pairs
	pairs, err := preparePairs(dB, users)
	if err != nil {
		return nil, err
	}
	// pick the best pairs
	return greedyMatcher(users, pairs)
}

func preparePairs(dB db.DB, users []model.User) ([]model.Match, error) {
	var pairs []model.Match
	for i := 0; i < len(users)-1; i++ {
		for j := i + 1; j < len(users); j++ {
			if users[i].ID == users[j].ID {
				continue
			}
			pairs = append(pairs, model.Match{A: users[i], B: users[j]})
		}
	}
	return dB.AddAndGetPairs(pairs)
}

// greedyMatcher implements the greedy algorithm for the minimum weight perfect matching problem
func greedyMatcher(users []model.User, pairs []model.Match) ([]model.Match, error) {
	// shuffle the pairs to remove bias towards lower ids
	swapPair := func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	}
	rand.Shuffle(len(pairs), swapPair)

	// stable sort the pairs by the least matches they have
	byFrequency := func(i, j int) bool {
		return pairs[i].Freq < pairs[j].Freq
	}
	sort.SliceStable(pairs, byFrequency)

	// number of matches is the ceil of the half of the number of users
	numMatches := len(users)/2 + len(users)%2

	matches := make([]model.Match, 0, numMatches)
	matched := make(map[int]bool, len(users))
	for _, user := range users {
		matched[user.ID] = false
	}

	// iterate over the pairs and match unmatched users
	for _, pair := range pairs {
		// break the loop if we have enough matches
		if len(matches) == numMatches-len(users)%2 {
			break
		}
		ok1 := matched[pair.A.ID]
		ok2 := matched[pair.B.ID]
		if ok1 || ok2 {
			continue
		}
		matches = append(matches, pair)
		matched[pair.A.ID] = true
		matched[pair.B.ID] = true
	}

	// in the case of an odd number of users, we need to add the last unmatched user
	if len(matches) == numMatches-1 {
		var unmatched int
		for id, found := range matched {
			if !found {
				unmatched = id
				break
			}
		}
		if unmatched == 0 {
			return nil, errors.New("no unmatched user found but one match is missing")
		}
		// match the unmatched user to the first available pair
		for _, pair := range pairs {
			if pair.A.ID == unmatched || pair.B.ID == unmatched {
				matches = append(matches, pair)
				break
			}
		}
	}

	if len(matches) != numMatches {
		return nil, fmt.Errorf("could not generate %d matches", numMatches)
	}

	return matches, nil
}
