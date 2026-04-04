package handler

import (
	"strconv"

	"message-push-system/internal/service"
	"message-push-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// MemberHandler 成员处理器
type MemberHandler struct {
	memberService service.MemberService
}

// NewMemberHandler 创建成员处理器实例
func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// AddMemberRequest 添加成员请求
type AddMemberRequest struct {
	Email string `json:"email" binding:"required"`
}

// Add 添加成员
func (h *MemberHandler) Add(c *gin.Context) {
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	member, err := h.memberService.Add(groupID, req.Email)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, member)
}

// Remove 删除成员
func (h *MemberHandler) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid member id")
		return
	}

	if err := h.memberService.Remove(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ListByGroup 列出群组的所有成员
func (h *MemberHandler) ListByGroup(c *gin.Context) {
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	members, err := h.memberService.ListByGroup(groupID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, members)
}
