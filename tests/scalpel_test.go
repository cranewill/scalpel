package tests

import (
	"fmt"
	"log"
	"scalpel"
	"testing"
)

type PlayerData struct {
	Uuid     string
	Sign     SignInfo
	SomeData Some
	Array    []int32
	Map      map[int32]int32
	MapSlice map[int32][]int32
	MapMap   map[int32]map[string]string
}

type Some struct {
	Id   int32
	Uuid string
	Pool map[int32]SignInfo
}

type SignInfo struct {
	Rewards []Reward
	Day     int32
}

type Reward struct {
	Category int32
	ConfId   int32
	Amount   int64
}

var player = PlayerData{
	Uuid: "player_1",
	Sign: SignInfo{
		Day: 1,
		Rewards: []Reward{
			{Category: 1, ConfId: 1, Amount: 100},
			{Category: 2, ConfId: 2, Amount: 1000},
		},
	},
	SomeData: Some{Id: 100, Uuid: "100", Pool: map[int32]SignInfo{
		1: {
			Rewards: []Reward{
				{Category: 1, ConfId: 1, Amount: 100},
				{Category: 2, ConfId: 2, Amount: 1000},
			},
			Day: 1,
		},
		2: {
			Rewards: []Reward{
				{Category: 11, ConfId: 11, Amount: 100},
				{Category: 22, ConfId: 22, Amount: 1000},
			},
			Day: 2,
		},
	}},
	Array: []int32{1, 2, 3},
	Map: map[int32]int32{
		1: 100,
		2: 200,
	},
	MapSlice: map[int32][]int32{
		1: {1, 2, 3},
	},
	MapMap: map[int32]map[string]string{
		1: {
			"apple":  "good",
			"orange": "nice",
		},
		2: {
			"dog":  "cute",
			"frog": "immortal",
		},
	},
}

func TestSetField(t *testing.T) {
	var err error
	params_uuid := []string{"Uuid"}
	params_sign_day := []string{"Sign", "Day"}
	params_sign_reward_confid := []string{"Sign", "Rewards", "0", "ConfId"}
	params_some_uuid := []string{"SomeData", "Uuid"}
	params_some_poll_reward_amount := []string{"SomeData", "Pool", "2", "Rewards", "1", "Amount"}
	params_map := []string{"Map", "2"}
	params_array := []string{"Array", "1"}
	params_map_slice := []string{"MapSlice", "1", "2"}
	params_map_map := []string{"MapMap", "2", "frog"}

	err = scalpel.SetField(&player, params_uuid, "new_playerId")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_sign_day, "7")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_sign_reward_confid, "12345")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_some_uuid, "some_uuid")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_some_poll_reward_amount, "777")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_map, "1000")
	if err != nil {
		panic(err)
	}
	err = scalpel.SetField(&player, params_array, "2222")
	if err != nil {
		panic(err)
	}
	err = scalpel.SetField(&player, params_map_slice, "3333")
	if err != nil {
		log.Println(err.Error())
	}
	err = scalpel.SetField(&player, params_map_map, "cute")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(fmt.Sprintf("%+v", player))
}

func TestSetFieldNotPtr(t *testing.T) {
	var err error
	params_uuid := []string{"Uuid"} // ???????????????????????????
	err = scalpel.SetField(player, params_uuid, "new_playerId")
	if err != nil {
		log.Println(err.Error())
	}
}

func TestSetFieldWrongPath(t *testing.T) {
	var err error
	params_uuid := []string{"Uuid", "Category"} // ??????????????????
	err = scalpel.SetField(&player, params_uuid, "1")
	if err != nil {
		log.Println(err.Error())
	}
}

func TestSetFieldWrongKind(t *testing.T) {
	var err error
	params1 := []string{"Sign", "Rewards", "good"} // ?????????key??????????????????
	err = scalpel.SetField(&player, params1, "1")
	if err != nil {
		log.Println(err.Error())
	}
	params2 := []string{"Sign", "Rewards", "0", "ConfId"}
	err = scalpel.SetField(&player, params2, "confId") // ?????????int32?????????????????????
	if err != nil {
		log.Println(err.Error())
	}
}

func TestSetFieldInvalidIndex(t *testing.T) {
	var err error
	params1 := []string{"Sign", "Rewards", "1000"} // ??????range
	err = scalpel.SetField(&player, params1, "1")
	if err != nil {
		log.Println(err.Error())
	}

	params2 := []string{"Sign", "Rewards", "-1"} // ??????0
	err = scalpel.SetField(&player, params2, "1")
	if err != nil {
		log.Println(err.Error())
	}
}
