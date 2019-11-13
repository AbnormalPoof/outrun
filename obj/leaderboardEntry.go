package obj

import "time"

type LeaderboardEntry struct {
    FriendID          string `json:"friendId"`
    Name              string `json:"name"`
    URL               string `json:"url"`
    Grade             int64  `json:"grade"`
    ExposeOnline      int64  `json:"exposeOnline"`
    RankingScore      int64  `json:"rankingScore"`
    RankChanged       int64  `json:"rankChanged"`
    IsSentEnergy      int64  `json:"energyFlg"`
    ExpireTime        int64  `json:"expireTime"`
    NumRank           int64  `json:"numRank"`
    LoginTime         int64  `json:"loginTime"`
    CharacterID       string `json:"charaId"`
    CharacterLevel    int64  `json:"characterLevel"`
    SubcharacterID    string `json:"subCharaId"`
    SubcharacterLevel int64  `json:"subCharaLevel"`
    MainChaoID        string `json:"mainChaoId"`
    MainChaoLevel     int64  `json:"mainChaoLevel"`
    SubChaoID         string `json:"subChaoId"`
    SubChaoLevel      int64  `json:"subChaoLevel"`
    Language          int64  `json:"language"`
    League            int64  `json:"league"`
    MaxScore          int64  `json:"maxScore"`
}

func NewLeaderboardEntry(fid, n, url string, g, eo, rs, rc, ise, et, nr, lt int64, cid string, cl int64, schid string, schl int64, mcid string, mcl int64, scid string, scl, lang, league, maxScore int64) LeaderboardEntry {
    return LeaderboardEntry{
        fid,
        n,
        url,
        g,
        eo,
        rs,
        rc,
        ise,
        et,
        nr,
        lt,
        cid,
        cl,
        schid,
        schl,
        mcid,
        mcl,
        scid,
        scl,
        lang,
        league,
        maxScore,
    }
}

func DefaultLeaderboardEntry(uid string) LeaderboardEntry {
    return NewLeaderboardEntry(
        uid,
        "",
        "",
        0,
        1,
        0,
        0,
        0,
        0,
        0,
        time.Now().Unix(), // this should be player.LastLogin!
        "0",
        0,
        "-1",
        0,
        "-1",
        0,
        "-1",
        0,
        1,
        0,
        0,
    )
}
