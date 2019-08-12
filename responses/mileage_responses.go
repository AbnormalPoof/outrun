package responses

import (
	"github.com/fluofoxxo/outrun/consts"
	"github.com/fluofoxxo/outrun/objects"
	"github.com/fluofoxxo/outrun/playerdata"
)

type MileageDataResponse struct {
	BaseResponse
	MileageFriendList     []objects.MileageFriend `json:"mileageFriendList"`
	pdata.MileageMapState `json:"mileageMapState"`
}

func NewMileageDataResponse(base BaseInfo, mileageFriendList []objects.MileageFriend, mileageMapState pdata.MileageMapState) MileageDataResponse {
	br := NewBaseResponse(base)
	mdr := MileageDataResponse{
		br,
		mileageFriendList,
		mileageMapState,
	}
	return mdr
}

func DefaultMileageDataResponse(base BaseInfo) MileageDataResponse {
	mileageMapState := pdata.MileageMapState{
		consts.MILE_EPISODE,
		consts.MILE_CHAPTER,
		consts.MILE_POINT,
		consts.MILE_MAPDISTANCE,
		consts.MILE_NUMBOSSATTACK,
		consts.MILE_STAGEDISTANCE,
		consts.MILE_STAGETOTALSCORE,
		consts.MILE_STAGEMAXSCORE,
		consts.MILE_CHAPTERSTARTTIME,
	}
	mdr := NewMileageDataResponse(
		base,
		[]objects.MileageFriend{},
		mileageMapState,
	)
	return mdr
}
