// Copyright 2018, gossiper project Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package guid

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// Web-safe chars ordered by ASCII.
	webSafeChars = "-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
)

var (
	// Timestamp of last push, used to prevent local collisions if you push
	// twice in one ms.
	lastPushTimeMs int64

	// Generate 72-bits of randomness which get turned into 12 characters and
	// appended to the timestamp to prevent collisions with other clients.
	// Store the last characters we generated because in the event of a
	// collision, we'll use those same characters except "incremented" by one.
	lastRandChars [12]int
	mu            sync.Mutex
	seed          *rand.Rand
)

func init() {
	// seed to get randomness
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func genRandPart() {
	for i := 0; i < len(lastRandChars); i++ {
		lastRandChars[i] = seed.Intn(64)
	}
}

// New creates a new random, unique id
func New() string {
	var id [8 + 12]byte
	mu.Lock()
	timeMs := time.Now().UTC().UnixNano() / 1e6
	if timeMs == lastPushTimeMs {
		// increment lastRandChars
		for i := 0; i < 12; i++ {
			lastRandChars[i]++
			if lastRandChars[i] < 64 {
				break
			}
			// increment the next byte
			lastRandChars[i] = 0
		}
	} else {
		genRandPart()
	}
	lastPushTimeMs = timeMs
	// put random as the second part
	for i := 0; i < 12; i++ {
		id[19-i] = webSafeChars[lastRandChars[i]]
	}
	mu.Unlock()

	// put current time at the beginning
	for i := 7; i >= 0; i-- {
		n := int(timeMs % 64)
		id[i] = webSafeChars[n]
		timeMs = timeMs / 64
	}
	return string(id[:])
}
