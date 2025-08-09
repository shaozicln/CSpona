package api

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"gorm.io/gorm"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

// 全局配置变量
var (
	cfg           *ini.File
	senderEmail   string
	smtpServer    string
	smtpPort      string
	checkInterval time.Duration
)

// 初始化配置
func init() {
	var err error
	// 加载配置文件
	cfg, err = ini.Load("config.ini")
	if err != nil {
		panic(fmt.Sprintf("无法加载配置文件: %v", err))
	}

	// 读取邮件配置
	senderEmail = cfg.Section("email").Key("SenderEmail").MustString("")
	smtpServer = cfg.Section("email").Key("SmtpServer").MustString("")
	smtpPort = cfg.Section("email").Key("SmtpPort").MustString("")

	// 读取检查间隔配置
	checkIntervalVal, err := cfg.Section("reminder").Key("check_interval").Int()
	if err != nil {
		panic(fmt.Sprintf("解析检查间隔失败: %v", err))
	}
	checkInterval = time.Duration(checkIntervalVal) * time.Minute

	// 验证配置
	if senderEmail == "" || smtpServer == "" || smtpPort == "" {
		panic("邮件配置不完整，请检查config.ini")
	}
}

type Reminder struct {
	Id        uint      `gorm:"primaryKey"`
	UserId    uint      `gorm:"not null;type:int unsigned" column:"user_id"`
	Title     string    `gorm:"type:varchar(255);not null" column:"title"`
	Content   string    `gorm:"type:text" column:"content"` // 对应前端的 description
	RemindAt  time.Time `gorm:"type:datetime;not null" column:"remind_at"`
	CreatedAt time.Time `gorm:"type:timestamp" column:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" column:"updated_at"`
	User      User      `gorm:"foreignKey:UserId"` // 关联用户表
}

func (Reminder) TableName() string {
	return "reminders"
}

// 创建提醒
func PostReminder(c *gin.Context) {
	var reminder Reminder
	// 获取当前登录用户ID（实际项目中需要从Token中解析）
	// 这里简化处理，实际应通过认证中间件获取
	userIdStr := c.PostForm("user_id")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 解析提醒时间
	remindAtStr := c.PostForm("remind_at")
	remindAt, err := time.Parse("2006-01-02 15:04:05", remindAtStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的时间格式，正确格式: 2006-01-02 15:04:05"})
		return
	}

	// 组装提醒数据
	reminder.UserId = uint(userId)
	reminder.Title = c.PostForm("title")
	reminder.Content = c.PostForm("content")
	reminder.RemindAt = remindAt
	reminder.CreatedAt = time.Now()
	reminder.UpdatedAt = time.Now()

	// 保存到数据库
	if err := db.Create(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建提醒失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提醒创建成功", "data": reminder})
}

// 获取用户的所有提醒
func GetUserReminders(c *gin.Context) {
	userIdStr := c.Param("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "无效的用户ID",
			"detail": err.Error(),
		})
		return
	}

	var reminders []Reminder
	result := db.Where("user_id = ?", int(userId)).
		Order("remind_at asc").
		Find(&reminders)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询提醒失败",
			"detail": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminders})
}

// 更新提醒
func PutReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	var reminder Reminder
	if err := db.First(&reminder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		return
	}

	// 更新标题
	title := c.PostForm("title")
	if title != "" {
		reminder.Title = title
	}

	// 更新内容
	content := c.PostForm("content")
	if content != "" {
		reminder.Content = content
	}

	// 更新提醒时间
	remindAtStr := c.PostForm("remind_at")
	if remindAtStr != "" {
		remindAt, err := time.Parse("2006-01-02 15:04:05", remindAtStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的时间格式"})
			return
		}
		reminder.RemindAt = remindAt
	}

	reminder.UpdatedAt = time.Now()

	if err := db.Save(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新提醒失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reminder})
}

// 删除提醒
func DeleteReminder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	var reminder Reminder
	if err := db.First(&reminder, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		return
	}

	if err := db.Delete(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除提醒失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提醒删除成功"})
}

// 带SSL的邮件发送函数，专门处理465端口
func sendMailWithSSL(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 创建TLS连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         smtpServer, // 服务器名称必须与证书一致
		InsecureSkipVerify: false,      // 生产环境不跳过证书验证
	})
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	// 进行认证
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %w", err)
		}
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	for _, recp := range to {
		if err := client.Rcpt(recp); err != nil {
			return fmt.Errorf("设置收件人失败(%s): %w", recp, err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送数据失败: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("发送邮件内容失败: %w", err)
	}

	// 完成发送
	if err := w.Close(); err != nil {
		return fmt.Errorf("关闭数据连接失败: %w", err)
	}

	return client.Quit()
}

// 发送邮件
func SendReminderEmail(c *gin.Context) {
	reminderIdStr := c.Param("id")
	reminderId, err := strconv.ParseUint(reminderIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提醒ID"})
		return
	}

	// 1. 查询提醒记录
	var reminder Reminder
	if err := db.First(&reminder, reminderId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "提醒不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询提醒失败: " + err.Error()})
		}
		return
	}

	// 2. 查询用户信息（获取收件人邮箱）
	var user User
	if err := db.First(&user, reminder.UserId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户信息失败: " + err.Error()})
		return
	}

	// 3. 收件人邮箱（用户邮箱）
	recipientEmail := user.Email
	if recipientEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未设置邮箱"})
		return
	}

	// 4. 获取SMTP授权码（解密）
	smtpPass := cfg.Section("email").Key("SmtpPass").MustString("")
	if smtpPass == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发件人SMTP授权码未配置"})
		return
	}

	// 5. 准备邮件内容
	subject := "【CSpona任务提醒】" + reminder.Title
	body := `
		<p>` + reminder.Content + `</p>
		<p style="text-align: right; margin-top: 20px;">From CSpona</p>
	`
	mailMsg := []byte(
		"From: " + senderEmail + "\r\n" +
			"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n\r\n" +
			body)

	// 6. 发送邮件
	addr := smtpServer + ":" + smtpPort
	auth := smtp.PlainAuth("", senderEmail, smtpPass, smtpServer)

	// 使用带SSL的发送函数（适合465端口）
	err = sendMailWithSSL(addr, auth, senderEmail, []string{recipientEmail}, mailMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "发送邮件失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "邮件提醒发送成功"})
}

// 定时检查并发送到期提醒（程序启动时初始化）
func InitReminderChecker() {
	// 使用从配置文件读取的检查间隔
	ticker := time.NewTicker(checkInterval)
	go func() {
		for range ticker.C {
			now := time.Now()
			// 查询当前时间前后1分钟内需要提醒的记录
			var reminders []Reminder
			db.Where("remind_at BETWEEN ? AND ?", now.Add(-1*time.Minute), now.Add(1*time.Minute)).
				Find(&reminders)

			for _, reminder := range reminders {
				// 调用发送邮件逻辑
				if err := sendReminderEmailToUser(reminder); err != nil {
					fmt.Printf("发送提醒邮件失败(提醒ID: %d): %v\n", reminder.Id, err)
				}
			}
		}
	}()
}

// 内部邮件发送函数
func sendReminderEmailToUser(reminder Reminder) error {
	var user User
	if err := db.First(&user, reminder.UserId).Error; err != nil {
		return fmt.Errorf("查询用户信息失败: %w", err)
	}

	// 收件人邮箱（用户邮箱）
	recipientEmail := user.Email
	if recipientEmail == "" {
		return errors.New("用户未设置邮箱")
	}

	// 获取SMTP授权码
	smtpPass := cfg.Section("email").Key("SmtpPass").MustString("")
	if smtpPass == "" {
		return errors.New("发件人SMTP授权码未配置")
	}

	// 准备邮件内容
	subject := "【CSpona任务提醒】" + reminder.Title
	body := `
		<p>` + reminder.Content + `</p>
		<p style="text-align: right; margin-top: 20px;">From CSpona</p>
	`
	mailMsg := []byte(
		"From: " + senderEmail + "\r\n" +
			"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n\r\n" +
			body)

	// 发送邮件
	addr := smtpServer + ":" + smtpPort
	auth := smtp.PlainAuth("", senderEmail, smtpPass, smtpServer)

	return sendMailWithSSL(addr, auth, senderEmail, []string{recipientEmail}, mailMsg)
}
