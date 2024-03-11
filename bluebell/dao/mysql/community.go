package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		zap.L().Warn("there is no community in db")
		err = nil
	}
	return
}

func CommunityDetailById(id int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := "select community_id, community_name, introduction, create_time from community where id = ?"
	err = db.Get(communityDetail, sqlStr, id)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
	}
	return
}
