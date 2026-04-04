package handler

import (
	"strconv"

	"message-push-system/internal/service"
	"message-push-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// LogHandler 推送日志处理器
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler 创建推送日志处理器实例
func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

// List 列出推送日志
func (h *LogHandler) List(c *gin.Context) {
	// 检查是否有 message_id 查询参数
	messageIDStr := c.Query("message_id")
	if messageIDStr != "" {
		messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid message_id")
			return
		}

		logs, err := h.logService.ListByMessage(messageID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, logs)
		return
	}

	// 列出所有日志
	logs, err := h.logService.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, logs)
}
