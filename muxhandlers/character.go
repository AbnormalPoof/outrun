package muxhandlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fluofoxxo/outrun/analytics"
	"github.com/fluofoxxo/outrun/analytics/factors"
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/db"
	"github.com/fluofoxxo/outrun/emess"
	"github.com/fluofoxxo/outrun/enums"
	"github.com/fluofoxxo/outrun/helper"
	"github.com/fluofoxxo/outrun/requests"
	"github.com/fluofoxxo/outrun/responses"
	"github.com/fluofoxxo/outrun/status"
)

func ChangeCharacter(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.ChangeCharacterRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	mainCharaID := request.MainCharaID
	subCharaID := request.SubCharaID
	if mainCharaID != "-1" {
		index := player.IndexOfChara(mainCharaID)
		if index == -1 {
			helper.Warn("Bad mainCharaID '%s'", mainCharaID)
			return
		}
		if player.CharacterState[index].Status != enums.CharacterStatusLocked {
			player.PlayerState.MainCharaID = mainCharaID
		} else {
			helper.Warn("Bad attempt to change main chara to '%s'", mainCharaID)
			baseInfo.StatusCode = status.DataMismatch
		}
		_, err = analytics.Store(player.ID, factors.AnalyticTypeChangeMainCharacter)
		if err != nil {
			helper.WarnErr("Error storing analytics (AnalyticTypeChangeMainCharacter)", err)
		}
	}
	if subCharaID != "-1" {
		index := player.IndexOfChara(subCharaID)
		if index == -1 {
			helper.Warn("Bad subCharaID '%s'", subCharaID)
			return
		}
		if player.CharacterState[index].Status != enums.CharacterStatusLocked {
			player.PlayerState.SubCharaID = subCharaID
		} else {
			helper.Warn("Bad attempt to change sub chara to '%s'", subCharaID)
			baseInfo.StatusCode = status.DataMismatch
		}
		_, err = analytics.Store(player.ID, factors.AnalyticTypeChangeSubCharacter)
		if err != nil {
			helper.WarnErr("Error storing analytics (AnalyticTypeChangeSubCharacter)", err)
		}
	}
	if player.PlayerState.MainCharaID == player.PlayerState.SubCharaID {
		helper.InvalidRequest()
		return
	}
	db.SavePlayer(player)

	response := responses.ChangeCharacter(baseInfo, player.PlayerState)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func UpgradeCharacter(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.UpgradeCharacterRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	charaID := request.CharacterID
	abilityID := request.AbilityID
	abilityIDFromStr, err := strconv.Atoi(abilityID)
	if err != nil {
		helper.Err("Error in strconv.Atoi", err)
		return
	}

	sum := func(in []int64) int64 {
		v := int64(0)
		for _, val := range in {
			v += val
		}
		return v
	}

	abilityIndex := abilityIDFromStr - 120000 // minus enums.UpgradeAbilityInvincibility
	index := player.IndexOfChara(charaID)
	abilitySum := sum(player.CharacterState[index].AbilityLevel)
	amountNeedToBePaid := player.CharacterState[index].Cost - player.CharacterState[index].Exp
	if player.PlayerState.NumRings-amountNeedToBePaid < 0 {
		baseInfo.StatusCode = status.NotEnoughRings
	} else {
		if abilitySum < 100 {
			// check if character is valid here
			levelIncrease, ok := consts.UpgradeIncreases[charaID]
			if !ok {
				helper.InternalErr("Error getting level increase", fmt.Errorf("key '%v' not found in consts.UpgradeIncreases", charaID))
				return
			}
			player.CharacterState[index].AbilityLevel[abilityIndex]++
			player.CharacterState[index].Level += 1
			player.CharacterState[index].Exp = 0 // reset exp
			player.CharacterState[index].Cost += levelIncrease
			player.PlayerState.NumRings -= amountNeedToBePaid
			if player.CharacterState[index].Level >= 100 {
				player.CharacterState[index].Status = enums.CharacterStatusMaxLevel
			}
			db.SavePlayer(player)
		} else {
			baseInfo.StatusCode = status.CharacterLevelLimit
		}
	}

	respPlayer := player
	if request.Version == "1.1.4" { // must send fewer characters
		// only get first 21 characters
		// TODO: enforce order 300000 to 300020?
		//cState = cState[:len(cState)-(len(cState)-10)]
		cState := respPlayer.CharacterState
		cState = cState[:16]
		helper.DebugOut("cState length: " + strconv.Itoa(len(cState)))
		helper.DebugOut("Sent character IDs: ")
		for _, char := range cState {
			helper.DebugOut(char.ID)
		}
		respPlayer.CharacterState = cState
	}
	response := responses.DefaultUpgradeCharacter(baseInfo, respPlayer)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func UnlockedCharacter(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.UnlockedCharacterRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}

	player, err := helper.GetCallingPlayer(true)
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Suspended {
		baseInfo.StatusCode = status.MissingPlayer
		err = helper.SendResponse(responses.NewBaseResponse(baseInfo))
		if err != nil {
			helper.InternalErr("Error sending response", err)
			return
		}
		return
	}

	characterToBuy := request.CharacterID
	charaIndex := player.IndexOfChara(characterToBuy)
	chara := player.CharacterState[charaIndex] // get actual character object
	buyWith := request.ItemID
	/*
		The reason that the lines about levels are commented is because
		the game acts as if the levels shouldn't have changed at all, which
		gives credence to the idea that the vanilla game wouldn't downlevel
		the characters.
		Looking back on it, this makes sense...
	*/
	helper.DebugOut("Pre:")
	if buyWith == enums.ItemIDStrRing { // is buying with rings
		ringCost := chara.Price
		if ringCost > player.PlayerState.NumRings { // cannot buy
			helper.DebugOut("Player can't pay with rings (Has %v)", player.PlayerState.NumRings)
			baseInfo.StatusCode = status.NotEnoughRings
		} else { // can buy with rings
			helper.DebugOut("NumRings: %v", player.PlayerState.NumRings)
			//helper.DebugOut(sp("CharacterState[%v].Level: %v", charaIndex, player.CharacterState[charaIndex].Level))
			helper.DebugOut("CharacterState[%v].Status: %v", charaIndex, player.CharacterState[charaIndex].Status)
			helper.DebugOut("CharacterState[%v].Star: %v", charaIndex, player.CharacterState[charaIndex].Star)
			player.PlayerState.NumRings -= ringCost
			//player.CharacterState[charaIndex].Level = 0
			if player.CharacterState[charaIndex].Status == enums.CharacterStatusUnlocked || player.CharacterState[charaIndex].Status == enums.CharacterStatusMaxLevel { // character already owned, so just limit break
				player.CharacterState[charaIndex].Star++
			} else if player.CharacterState[charaIndex].Status == enums.CharacterStatusLocked { // character not already owned, so purchase them
				player.CharacterState[charaIndex].Status = enums.CharacterStatusUnlocked
			}
			db.SavePlayer(player)
		}
	} else if buyWith == enums.ItemIDStrRedRing { // is buying with red rings
		redRingCost := chara.PriceRedRings
		if redRingCost > player.PlayerState.NumRedRings { // cannot buy with red rings
			helper.DebugOut("Player can't pay with red rings (Has %v)", player.PlayerState.NumRedRings)
			baseInfo.StatusCode = status.NotEnoughRedRings
		} else { // can buy with red rings
			helper.DebugOut("NumRedRings: %v", player.PlayerState.NumRedRings)
			//helper.DebugOut(sp("CharacterState[%v].Level: %v", charaIndex, player.CharacterState[charaIndex].Level))
			helper.DebugOut("CharacterState[%v].Status: %v", charaIndex, player.CharacterState[charaIndex].Status)
			helper.DebugOut("CharacterState[%v].Star: %v", charaIndex, player.CharacterState[charaIndex].Star)
			player.PlayerState.NumRedRings -= redRingCost
			//player.CharacterState[charaIndex].Level = 0
			if player.CharacterState[charaIndex].Status == enums.CharacterStatusUnlocked || player.CharacterState[charaIndex].Status == enums.CharacterStatusMaxLevel { // character already owned, so just limit break
				player.CharacterState[charaIndex].Star++
			} else if player.CharacterState[charaIndex].Status == enums.CharacterStatusLocked { // character not already owned, so purchase them
				player.CharacterState[charaIndex].Status = enums.CharacterStatusUnlocked
			}
			db.SavePlayer(player)
		}
	} else { // didn't buy using rings or red rings...
		helper.Warn(fmt.Sprintf("Player '%s' (%v) tried to purchase a character without Rings or Red Rings!", player.Username, player.ID))
		baseInfo.StatusCode = status.OtherError
	}
	helper.DebugOut("Post:")
	helper.DebugOut("NumRings: %v", player.PlayerState.NumRings)
	helper.DebugOut("NumRedRings: %v", player.PlayerState.NumRedRings)
	helper.DebugOut("CharacterState[%v].Level: %v", charaIndex, player.CharacterState[charaIndex].Level)
	helper.DebugOut("CharacterState[%v].Status: %v", charaIndex, player.CharacterState[charaIndex].Status)
	helper.DebugOut("CharacterState[%v].Star: %v", charaIndex, player.CharacterState[charaIndex].Star)

	response := responses.DefaultUpgradeCharacter(baseInfo, player)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
