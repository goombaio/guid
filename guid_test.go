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

package guid_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/goombaio/guid"
)

const (
	// for manual debugging
	showIds = false
)

func Test(t *testing.T) {
	id1 := guid.New()
	if len(id1) != 20 {
		t.Fatalf("len(id1) != 20 (=%d)", len(id1))
	}
	id2 := guid.New()
	if len(id2) != 20 {
		t.Fatalf("len(id1) != 20 (=%d)", len(id2))
	}
	if id1 == id2 {
		t.Fatalf("generated same ids (id1: '%s', id2: '%s')", id1, id2)
	}
	if showIds {
		fmt.Printf("%s\n", id1)
		fmt.Printf("%s\n", id2)
	}
}

// As t.Fatalf() is not goroutine safe, use this closure.
func fail(t *testing.T, template string, args ...interface{}) {
	fmt.Printf(template, args...)
	fmt.Println()
	t.Fail()
}

func doMany(t *testing.T, wg *sync.WaitGroup) {
	ids := make(map[string]bool)
	prev := ""
	for i := 0; i < 1000000; i++ {
		id := guid.New()
		if _, exists := ids[id]; exists {
			fail(t, "generated duplicate id '%s'", id)
		}
		ids[id] = true
		if prev != "" {
			if id <= prev {
				fail(t, "id ('%s') must be > prev ('%s')", id, prev)
			}
		}
		prev = id
	}
	wg.Done()
}

func TestMany(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go doMany(t, &wg)
	}
	wg.Wait()
}
