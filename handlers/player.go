package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/fluofoxxo/outrun/consts"
    "github.com/fluofoxxo/outrun/cryption"
    "github.com/fluofoxxo/outrun/db"
    "github.com/fluofoxxo/outrun/helper"
    "github.com/fluofoxxo/outrun/requests"
    "github.com/fluofoxxo/outrun/responses"
)

func SetUserNameHandler(w http.ResponseWriter, r *http.Request) {
    recv := cryption.GetReceivedMessage(r)

    var request requests.SetUsernameRequest
    err := json.Unmarshal(recv, &request)
    if err != nil {
        log.Println("[ERR] (SetUserNameHandler) Error unmarshalling: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid request"))
        return
    }

    player, err := db.SessionIDToPlayer(request.SessionID)
    if err != nil {
        log.Println("[ERR] (SetUserNameHandler) Error getting player: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    player.Username = request.Username
    err = db.SavePlayer(player)
    if err != nil {
        log.Println("[ERR] (SetUserNameHandler) Error saving player: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    baseInfo := responses.NewBaseInfo(consts.EM_OK, 0, 0) // TODO: error handling in case the username is already taken!!
    resp := responses.NewBaseResponse(baseInfo)
    respJ, err := responses.ToJSON(resp)
    if err != nil {
        log.Println("[ERR] (SetUserNameHandler) Error marshalling: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    log.Println("[OUT] (SetUserNameHandler) All OK")
    helper.Respond([]byte(respJ), w)
}

func GetPlayerStateHandler(w http.ResponseWriter, r *http.Request) {
    request, err := helper.GetBasicRequest(r)
    if err != nil {
        log.Println("[ERR] (GetPlayerStateHandler) JSON unmarshal error: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid request"))
        return
    }

    baseInfo := responses.NewBaseInfo(consts.EM_OK, 0, 0)

    uid, err := db.SessionIDToUID(request.SessionID)
    if err != nil {
        log.Println("[ERR] (GetPlayerStateHandler) Error in resolving SID to UID: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }
    player, err := db.GetPlayerByUID(uid)
    if err != nil {
        log.Println("[ERR] (GetPlayerStateHandler) Error getting player: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    resp := responses.NewPlayerStateResponse(baseInfo, player.PlayerState)
    respJ, err := responses.ToJSON(resp)
    if err != nil {
        log.Println("[ERR] (GetPlayerStateHandler) Error in JSON marshalling: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }
    log.Println("[OUT] (GetPlayerStateHandler) All OK")
    helper.Respond([]byte(respJ), w)
}

func GetCharacterStateHandler(w http.ResponseWriter, r *http.Request) {
    request, err := helper.GetBasicRequest(r)
    if err != nil {
        log.Println("[ERR] (GetCharacterStateHandler) JSON unmarshal error: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid request"))
        return
    }
    baseInfo := responses.NewBaseInfo(consts.EM_OK, 0, 0)
    uid, err := db.SessionIDToUID(request.SessionID)
    if err != nil {
        log.Println("[ERR] (GetCharacterStateHandler) Error in resolving SID to UID: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }
    player, err := db.GetPlayerByUID(uid)
    if err != nil {
        log.Println("[ERR] (GetCharacterStateHandler) Error getting player: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    resp := responses.NewCharacterStateResponse(baseInfo, player.CharacterStates)
    respJ, err := responses.ToJSON(resp)
    if err != nil {
        log.Println("[ERR] (GetCharacterStateHandler) Error in JSON marshalling: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    log.Println("[OUT] (GetCharacterStateHandler) All OK")
    helper.Respond([]byte(respJ), w)
}

func GetChaoStateHandler(w http.ResponseWriter, r *http.Request) {
    request, err := helper.GetBasicRequest(r)
    if err != nil {
        log.Println("[ERR] (GetChaoStateHandler) JSON unmarshal error: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid request"))
        return
    }
    baseInfo := responses.NewBaseInfo(consts.EM_OK, 0, 0)
    uid, err := db.SessionIDToUID(request.SessionID)
    if err != nil {
        log.Println("[ERR] (GetChaoStateHandler) Error in resolving SID to UID: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }
    player, err := db.GetPlayerByUID(uid)
    if err != nil {
        log.Println("[ERR] (GetChaoStateHandler) Error getting player: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    resp := responses.NewChaoStateResponse(baseInfo, player.ChaoState)
    respJ, err := responses.ToJSON(resp)
    if err != nil {
        log.Println("[ERR] (GetChaoStateHandler) Error in JSON marshalling: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    log.Println("[OUT] (GetChaoStateHandler) All OK")
    helper.Respond([]byte(respJ), w)
}
