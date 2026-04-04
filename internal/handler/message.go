package handler

import (
	"message-push-system/internal/service"
	"message-push-system/internal/worker"
	"message-push-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageService service.MessageService
	pushWorker     *worker.PushWorker
}

// NewMessageHandler 创建消息处理器实例
func NewMessageHandler(messageService service.MessageService, pushWorker *worker.PushWorker) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		pushWorker:     pushWorker,
	}
}

// Create 创建消息并触发推送
func (h *MessageHandler) Create(c *gin.Context) {
	var req service.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	message, err := h.messageService.Create(&req)
	if err != nil {
		if err.Error() == "group not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// 提交推送任务
	h.pushWorker.Submit(&worker.PushTask{
		MessageID: message.ID,
		GroupID:   message.GroupID,
	})

	response.Success(c, message)
}
