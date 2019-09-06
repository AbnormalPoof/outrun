package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fluofoxxo/outrun/config"
	"github.com/fluofoxxo/outrun/cryption"
	"github.com/fluofoxxo/outrun/muxhandlers"
	"github.com/fluofoxxo/outrun/muxhandlers/muxobj"
)

const UNKNOWN_REQUEST_DIRECTORY = "logging/unknown_requests/"

var (
	LogExecutionTime = true
)

func OutputUnknownRequest(w http.ResponseWriter, r *http.Request) {
	recv := cryption.GetReceivedMessage(r)
	// make a new logging path
	timeStr := strconv.Itoa(int(time.Now().Unix()))
	os.MkdirAll(UNKNOWN_REQUEST_DIRECTORY, 0644)
	normalizedReq := strings.ReplaceAll(r.URL.Path, "/", "-")
	path := UNKNOWN_REQUEST_DIRECTORY + normalizedReq + "_" + timeStr + ".txt"
	err := ioutil.WriteFile(path, recv, 0644)
	if err != nil {
		log.Println("[OUT] UNABLE TO WRITE UNKNOWN REQUEST: " + err.Error())
		w.Write([]byte(""))
		return
	}
	log.Println("[OUT] !!!!!!!!!!!! Unknown request, output to " + path)
	w.Write([]byte(""))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := config.Parse("config.json")
	if err != nil {
		log.Printf("[INFO] No config file (config.json) found (%s), using defaults\n", err)
	} else {
		log.Println("[INFO] Config file (config.json) loaded")
	}

	h := muxobj.Handle
	mux := http.NewServeMux()
	LogExecutionTime = config.CFile.DoTimeLogging
	// Login
	mux.HandleFunc("/Login/login/", h(muxhandlers.Login, LogExecutionTime))
	mux.HandleFunc("/Sgn/sendApollo/", h(muxhandlers.SendApollo, LogExecutionTime))
	mux.HandleFunc("/Login/getVariousParameter/", h(muxhandlers.GetVariousParameter, LogExecutionTime))
	mux.HandleFunc("/Player/getPlayerState/", h(muxhandlers.GetPlayerState, LogExecutionTime))
	mux.HandleFunc("/Player/getCharacterState/", h(muxhandlers.GetCharacterState, LogExecutionTime))
	mux.HandleFunc("/Player/getChaoState/", h(muxhandlers.GetChaoState, LogExecutionTime))
	mux.HandleFunc("/Spin/getWheelOptions/", h(muxhandlers.GetWheelOptions, LogExecutionTime))
	mux.HandleFunc("/Game/getDailyChalData/", h(muxhandlers.GetDailyChallengeData, LogExecutionTime))
	mux.HandleFunc("/Message/getMessageList/", h(muxhandlers.GetMessageList, LogExecutionTime))
	mux.HandleFunc("/Store/getRedstarExchangeList/", h(muxhandlers.GetRedStarExchangeList, LogExecutionTime))
	mux.HandleFunc("/Game/getCostList/", h(muxhandlers.GetCostList, LogExecutionTime))
	mux.HandleFunc("/Event/getEventList/", h(muxhandlers.GetEventList, LogExecutionTime))
	mux.HandleFunc("/Game/getMileageData/", h(muxhandlers.GetMileageData, LogExecutionTime))
	mux.HandleFunc("/Game/getCampaignList/", h(muxhandlers.GetCampaignList, LogExecutionTime))
	mux.HandleFunc("/Chao/getChaoWheelOptions/", h(muxhandlers.GetChaoWheelOptions, LogExecutionTime))
	mux.HandleFunc("/Chao/getPrizeChaoWheelSpin/", h(muxhandlers.GetPrizeChaoWheelSpin, LogExecutionTime))
	mux.HandleFunc("/login/getInformation/", h(muxhandlers.GetInformation, LogExecutionTime))
	mux.HandleFunc("/Leaderboard/getWeeklyLeaderboardOptions/", h(muxhandlers.GetWeeklyLeaderboardOptions, LogExecutionTime))
	mux.HandleFunc("/Leaderboard/getLeagueData/", h(muxhandlers.GetLeagueData, LogExecutionTime))
	mux.HandleFunc("/Leaderboard/getWeeklyLeaderboardEntries/", h(muxhandlers.GetWeeklyLeaderboardEntries, LogExecutionTime))
	mux.HandleFunc("/Player/setUserName/", h(muxhandlers.SetUsername, LogExecutionTime))
	mux.HandleFunc("/login/getTicker/", h(muxhandlers.GetTicker, LogExecutionTime))
	mux.HandleFunc("/Login/loginBonus/", h(muxhandlers.LoginBonus, LogExecutionTime))
	// Timed mode
	mux.HandleFunc("/Game/quickActStart/", h(muxhandlers.QuickActStart, LogExecutionTime))
	mux.HandleFunc("/Game/quickPostGameResults/", h(muxhandlers.QuickPostGameResults, LogExecutionTime))
	// Story mode
	mux.HandleFunc("/Game/actStart/", h(muxhandlers.ActStart, LogExecutionTime))
	// Retry
	mux.HandleFunc("/Game/actRetry/", h(muxhandlers.ActRetry, LogExecutionTime))
	// Gameplay
	mux.HandleFunc("/Game/getFreeItemList/", h(muxhandlers.GetFreeItemList, LogExecutionTime))
	mux.HandleFunc("/Game/postGameResults/", h(muxhandlers.PostGameResults, LogExecutionTime))
	// Misc.
	mux.HandleFunc("/Character/changeCharacter/", h(muxhandlers.ChangeCharacter, LogExecutionTime))
	mux.HandleFunc("/Character/upgradeCharacter/", h(muxhandlers.UpgradeCharacter, LogExecutionTime))
	mux.HandleFunc("/Chao/equipChao/", h(muxhandlers.EquipChao, LogExecutionTime))
	// Shop
	mux.HandleFunc("/Store/redstarExchange/", h(muxhandlers.RedStarExchange, LogExecutionTime))

	if config.CFile.LogUnknownRequests {
		mux.HandleFunc("/", OutputUnknownRequest)
	}

	port := config.CFile.Port
	log.Printf("Starting server on port %s\n", port)
	panic(http.ListenAndServe(":"+port, mux))
}
