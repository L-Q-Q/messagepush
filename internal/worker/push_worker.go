package worker

import (
	"fmt"
	"log"
	"time"

	"message-push-system/internal/model"
	"message-push-system/internal/repository"
	"message-push-system/internal/smtp"
)

// PushWorker 推送工作器
type PushWorker struct {
	taskChan    chan *PushTask
	smtpClient  *smtp.SMTPClient
	messageRepo repository.MessageRepository
	memberRepo  repository.MemberRepository
	logRepo     repository.LogRepository
	groupRepo   repository.GroupRepository
	stopChan    chan struct{}
}

// PushTask 推送任务
type PushTask struct {
	MessageID int64
	GroupID   int64
}

// NewPushWorker 创建推送工作器实例
func NewPushWorker(
	smtpClient *smtp.SMTPClient,
	messageRepo repository.MessageRepository,
	memberRepo repository.MemberRepository,
	logRepo repository.LogRepository,
	groupRepo repository.GroupRepository,
) *PushWorker {
	return &PushWorker{
		taskChan:    make(chan *PushTask, 100),
		smtpClient:  smtpClient,
		messageRepo: messageRepo,
		memberRepo:  memberRepo,
		logRepo:     logRepo,
		groupRepo:   groupRepo,
		stopChan:    make(chan struct{}),
	}
}

// Start 启动工作器
func (w *PushWorker) Start() {
	go func() {
		log.Println("Push worker started")
		for {
			select {
			case task := <-w.taskChan:
				w.process(task)
			case <-w.stopChan:
				log.Println("Push worker stopped")
				return
			}
		}
	}()
}

// Stop 停止工作器
func (w *PushWorker) Stop() {
	close(w.stopChan)
}

// Submit 提交推送任务
func (w *PushWorker) Submit(task *PushTask) {
	w.taskChan <- task
}

// process 处理推送任务
func (w *PushWorker) process(task *PushTask) {
	log.Printf("Processing push task: message_id=%d, group_id=%d", task.MessageID, task.GroupID)

	// 更新消息状态为 processing
	if err := w.messageRepo.UpdateStatus(task.MessageID, model.StatusProcessing); err != nil {
		log.Printf("Failed to update message status to processing: %v", err)
		return
	}

	// 获取消息
	message, err := w.messageRepo.GetByID(task.MessageID)
	if err != nil {
		log.Printf("Failed to get message: %v", err)
		return
	}

	// 获取群组
	group, err := w.groupRepo.GetByID(task.GroupID)
	if err != nil {
		log.Printf("Failed to get group: %v", err)
		w.messageRepo.UpdateStatus(task.MessageID, model.StatusFailed)
		return
	}

	// 获取群组成员
	members, err := w.memberRepo.ListByGroup(task.GroupID)
	if err != nil {
		log.Printf("Failed to get group members: %v", err)
		w.messageRepo.UpdateStatus(task.MessageID, model.StatusFailed)
		return
	}

	if len(members) == 0 {
		log.Printf("No members in group %d", task.GroupID)
		w.messageRepo.UpdateStatus(task.MessageID, model.StatusSuccess)
		return
	}

	// 构建 SMTP 配置
	emailConfig := &smtp.EmailConfig{
		Server:   group.SMTPServer,
		Port:     group.SMTPPort,
		Username: group.SMTPUsername,
		Password: group.SMTPPassword,
	}

	// 发送邮件给每个成员
	hasError := false
	for _, member := range members {
		email := &smtp.Email{
			From:    group.SMTPUsername,
			To:      member.Email,
			Subject: message.Subject,
			Body:    message.Body,
		}

		// 发送邮件
		err := w.smtpClient.Send(emailConfig, email)

		// 创建推送日志
		pushLog := &model.PushLog{
			MessageID: task.MessageID,
			Recipient: member.Email,
			CreatedAt: time.Now(),
		}

		if err != nil {
			log.Printf("Failed to send email to %s: %v", member.Email, err)
			pushLog.Status = model.StatusFailed
			pushLog.ErrorMessage = fmt.Sprintf("SMTP error: %v", err)
			hasError = true
		} else {
			log.Printf("Email sent successfully to %s", member.Email)
			pushLog.Status = model.StatusSuccess
		}

		// 保存推送日志
		if err := w.logRepo.Create(pushLog); err != nil {
			log.Printf("Failed to create push log: %v", err)
		}
	}

	// 更新消息状态
	finalStatus := model.StatusSuccess
	if hasError {
		finalStatus = model.StatusFailed
	}

	if err := w.messageRepo.UpdateStatus(task.MessageID, finalStatus); err != nil {
		log.Printf("Failed to update message status to %s: %v", finalStatus, err)
	}

	log.Printf("Push task completed: message_id=%d, status=%s", task.MessageID, finalStatus)
}
