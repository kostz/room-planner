package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestRoomPlanner_BigRoom(t *testing.T) {
	is := is.New(t)
	rp := NewRoomPlanner("tests/room-big.txt")
	rp.traverse()

	is.Equal(rp.rooms, map[string]RoomChairs{
		"_total": {
			"W": 14,
			"P": 7,
			"S": 3,
			"C": 1,
		},
		"balcony": {
			"W": 0,
			"P": 2,
			"S": 0,
			"C": 0,
		},
		"bathroom": {
			"W": 0,
			"P": 1,
			"S": 0,
			"C": 0,
		},
		"closet": {
			"W": 0,
			"P": 3,
			"S": 0,
			"C": 0,
		},
		"kitchen": {
			"W": 4,
			"P": 0,
			"S": 0,
			"C": 0,
		},
		"living room": {
			"W": 7,
			"P": 0,
			"S": 2,
			"C": 0,
		},
		"office": {
			"W": 2,
			"P": 1,
			"S": 0,
			"C": 0,
		},
		"sleeping room": {
			"W": 1,
			"P": 0,
			"S": 1,
			"C": 0,
		},
		"toilet": {
			"W": 0,
			"P": 0,
			"S": 0,
			"C": 1,
		},
	})
}

func TestRoomPlanner_SmallRoom(t *testing.T) {
	is := is.New(t)
	rp := NewRoomPlanner("tests/room-small.txt")
	rp.traverse()

	is.Equal(rp.rooms, map[string]RoomChairs{
		"_total": {
			"W": 0,
			"P": 8,
			"S": 0,
			"C": 0,
		},
		"closet": {
			"W": 0,
			"P": 8,
			"S": 0,
			"C": 0,
		},
	})
}
