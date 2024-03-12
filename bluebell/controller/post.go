package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数并校验参数
	post := new(models.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		// 校验参数
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
	}
	// 获取当前登录的用户
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorID = userID
	// 2. 插入数据库
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从URL中获取帖子的id）
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
	}
	// 2. 根据id查询帖子数据
	data, err := logic.GetPostById(id)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表处理函数
func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList(page, size) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	ResponseSuccess(c, data)
}
