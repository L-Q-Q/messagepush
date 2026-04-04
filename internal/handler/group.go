package handler

import (
	"strconv"

	"message-push-system/internal/service"
	"message-push-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// GroupHandler 群组处理器
type GroupHandler struct {
	groupService service.GroupService
}

// NewGroupHandler 创建群组处理器实例
func NewGroupHandler(groupService service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// Create 创建群组
func (h *GroupHandler) Create(c *gin.Context) {
	var req service.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	group, err := h.groupService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, group)
}

// Delete 删除群组
func (h *GroupHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.groupService.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// List 列出所有群组
func (h *GroupHandler) List(c *gin.Context) {
	groups, err := h.groupService.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, groups)
}
