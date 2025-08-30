// Package sdk provides the MiloraBot SDK interface
package sdk

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"time"

	"github.com/gorilla/websocket"

	"github.com/iamlibie/milonra-go/bot"
	"github.com/iamlibie/milonra-go/event"
	mplugin "github.com/iamlibie/milonra-go/plugin"
)

// Bot 接口，包装核心plugin.Bot接口
type Bot interface {
	WriteJSON(v interface{}) error // 线程安全地写入JSON数据
	GetSelfID() int64              // 返回机器人ID
}

// PluginFunc 插件函数类型：输入bot实例和消息事件，输出回复
type PluginFunc func(bot Bot, event *event.MessageEvent) string

// botAdapter 适配器，将sdk.Bot转换为plugin.Bot
type botAdapter struct {
	inner Bot
}

func (ba *botAdapter) WriteJSON(v interface{}) error {
	return ba.inner.WriteJSON(v)
}

func (ba *botAdapter) GetSelfID() int64 {
	return ba.inner.GetSelfID()
}

// MiloraBotConfig SDK配置结构
type MiloraBotConfig struct {
	// 服务配置
	Port         string        `json:"port"`          // 监听端口，默认 ":8080"
	Host         string        `json:"host"`          // 监听主机，默认 ""（所有接口）
	ReadTimeout  time.Duration `json:"read_timeout"`  // 读取超时，默认 15秒
	WriteTimeout time.Duration `json:"write_timeout"` // 写入超时，默认 15秒

	// 机器人配置
	BotID int64 `json:"bot_id"` // 机器人QQ号

	// WebSocket配置
	CheckOrigin func(*http.Request) bool `json:"-"` // 跨域检查函数

	// 日志配置
	EnableLog bool   `json:"enable_log"` // 是否启用日志，默认 true
	LogLevel  string `json:"log_level"`  // 日志级别：debug, info, warn, error

	// 插件配置
	PluginDir         string   `json:"plugin_dir"`          // 插件目录，默认 "./plugins"
	EnabledPlugins    []string `json:"enabled_plugins"`     // 启用的插件列表，空表示全部启用
	AutoLoadPlugins   bool     `json:"auto_load_plugins"`   // 是否自动加载插件目录中的插件
	PluginFilePattern string   `json:"plugin_file_pattern"` // 插件文件匹配模式，默认 "*.so"
}

// MiloraBot SDK主结构
type MiloraBot struct {
	config   *MiloraBotConfig
	server   *http.Server
	upgrader websocket.Upgrader
	bot      *bot.Bot
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewMiloraBot 创建新的MiloraBot实例
func NewMiloraBot(config *MiloraBotConfig) *MiloraBot {
	// 设置默认配置
	if config == nil {
		config = &MiloraBotConfig{}
	}

	if config.Port == "" {
		config.Port = ":8080"
	}

	if config.ReadTimeout == 0 {
		config.ReadTimeout = 15 * time.Second
	}

	if config.WriteTimeout == 0 {
		config.WriteTimeout = 15 * time.Second
	}

	if config.CheckOrigin == nil {
		config.CheckOrigin = func(r *http.Request) bool { return true }
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if config.PluginDir == "" {
		config.PluginDir = "./plugins"
	}

	if config.PluginFilePattern == "" {
		config.PluginFilePattern = "*.so"
	}

	// 默认启用自动加载插件
	config.AutoLoadPlugins = true

	// 默认启用日志
	config.EnableLog = true

	ctx, cancel := context.WithCancel(context.Background())

	mb := &MiloraBot{
		config: config,
		upgrader: websocket.Upgrader{
			CheckOrigin: config.CheckOrigin,
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// 创建HTTP服务器
	mb.server = &http.Server{
		Addr:         config.Port,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return mb
}

// RegisterPlugin 注册插件
func (mb *MiloraBot) RegisterPlugin(name string, pluginFunc PluginFunc) {
	// 将SDK插件包装为核心plugin包的插件
	wrappedFunc := func(bot mplugin.Bot, e *event.MessageEvent) string {
		return pluginFunc(&botAdapter{bot}, e)
	}
	mplugin.Register(name, wrappedFunc)
	if mb.config.EnableLog {
		log.Printf("🔌 插件已注册: %s", name)
	}
}

// SetPluginDir 设置插件目录
func (mb *MiloraBot) SetPluginDir(dir string) {
	mb.config.PluginDir = dir
}

// SetAutoLoadPlugins 设置是否自动加载插件
func (mb *MiloraBot) SetAutoLoadPlugins(enable bool) {
	mb.config.AutoLoadPlugins = enable
}

// LoadPluginsFromDir 从指定目录加载插件
func (mb *MiloraBot) LoadPluginsFromDir(dir string) error {
	if dir == "" {
		dir = mb.config.PluginDir
	}

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mb.config.EnableLog {
			log.Printf("📁 插件目录不存在，跳过自动加载: %s", dir)
		}
		return nil
	}

	// 查找插件文件
	pattern := filepath.Join(dir, mb.config.PluginFilePattern)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("查找插件文件失败: %v", err)
	}

	loadedCount := 0
	for _, file := range files {
		if err := mb.loadPluginFile(file); err != nil {
			if mb.config.EnableLog {
				log.Printf("❌ 加载插件文件失败 %s: %v", file, err)
			}
			continue
		}
		loadedCount++
	}

	if mb.config.EnableLog && loadedCount > 0 {
		log.Printf("🔌 从目录 %s 加载了 %d 个插件", dir, loadedCount)
	}

	return nil
}

// loadPluginFile 加载单个插件文件
func (mb *MiloraBot) loadPluginFile(filename string) error {
	p, err := plugin.Open(filename)
	if err != nil {
		return err
	}

	// 查找插件初始化函数
	initFunc, err := p.Lookup("Init")
	if err != nil {
		return fmt.Errorf("插件 %s 缺少 Init 函数: %v", filename, err)
	}

	// 调用初始化函数
	if initFn, ok := initFunc.(func()); ok {
		initFn()
		if mb.config.EnableLog {
			log.Printf("✅ 插件加载成功: %s", filepath.Base(filename))
		}
	} else {
		return fmt.Errorf("插件 %s 的 Init 函数签名不正确", filename)
	}

	return nil
}

// SetBotID 设置机器人ID
func (mb *MiloraBot) SetBotID(botID int64) {
	mb.config.BotID = botID
}

// SetPort 设置监听端口
func (mb *MiloraBot) SetPort(port string) {
	mb.config.Port = port
	if mb.server != nil {
		mb.server.Addr = port
	}
}

// SetLogLevel 设置日志级别
func (mb *MiloraBot) SetLogLevel(level string) {
	mb.config.LogLevel = level
}

// EnableLogging 启用或禁用日志
func (mb *MiloraBot) EnableLogging(enable bool) {
	mb.config.EnableLog = enable
}

// SetReadTimeout 设置读取超时
func (mb *MiloraBot) SetReadTimeout(timeout time.Duration) {
	mb.config.ReadTimeout = timeout
	if mb.server != nil {
		mb.server.ReadTimeout = timeout
	}
}

// SetWriteTimeout 设置写入超时
func (mb *MiloraBot) SetWriteTimeout(timeout time.Duration) {
	mb.config.WriteTimeout = timeout
	if mb.server != nil {
		mb.server.WriteTimeout = timeout
	}
}

// GetConfig 获取当前配置
func (mb *MiloraBot) GetConfig() *MiloraBotConfig {
	return mb.config
}

// handleWebSocket 处理WebSocket连接
func (mb *MiloraBot) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := mb.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if mb.config.EnableLog {
			log.Printf("WebSocket升级失败: %v", err)
		}
		return
	}
	defer conn.Close()

	if mb.config.EnableLog {
		log.Println("OneBot客户端已连接！")
	}

	// 创建机器人实例
	mb.bot = &bot.Bot{
		Conn:   conn,
		SelfID: mb.config.BotID,
	}

	// 消息处理循环
	for {
		select {
		case <-mb.ctx.Done():
			if mb.config.EnableLog {
				log.Println("收到停止信号，关闭WebSocket连接")
			}
			return
		default:
			var data map[string]interface{}
			err := conn.ReadJSON(&data)
			if err != nil {
				if mb.config.EnableLog {
					log.Printf("读取消息失败: %v", err)
				}
				return
			}

			// 交给Bot处理
			mb.bot.HandleMessage(data)
		}
	}
}

// Start 启动MiloraBot服务
func (mb *MiloraBot) Start() error {
	// 自动加载插件
	if mb.config.AutoLoadPlugins {
		if err := mb.LoadPluginsFromDir(""); err != nil {
			if mb.config.EnableLog {
				log.Printf("⚠️ 自动加载插件失败: %v", err)
			}
		}
	}

	// 设置路由
	http.HandleFunc("/", mb.handleWebSocket)

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 状态信息端点
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		plugins := mplugin.GetPlugins()
		status := map[string]interface{}{
			"status":     "running",
			"bot_id":     mb.config.BotID,
			"port":       mb.config.Port,
			"plugins":    len(plugins),
			"plugin_dir": mb.config.PluginDir,
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"%s","bot_id":%d,"port":"%s","plugins":%d}`,
			status["status"], status["bot_id"], status["port"], status["plugins"])
	})

	if mb.config.EnableLog {
		plugins := mplugin.GetPlugins()
		log.Printf("MiloraBot 启动中...")
		log.Printf("监听端口: %s", mb.config.Port)
		log.Printf("机器人ID: %d", mb.config.BotID)
		log.Printf("插件目录: %s", mb.config.PluginDir)
		log.Printf("已加载插件: %d个", len(plugins))
		log.Printf("健康检查: http://localhost%s/health", mb.config.Port)
		log.Printf("状态信息: http://localhost%s/status", mb.config.Port)
	}

	return mb.server.ListenAndServe()
}

// Stop 停止MiloraBot服务
func (mb *MiloraBot) Stop(timeout time.Duration) error {
	if mb.config.EnableLog {
		log.Println("正在停止MiloraBot服务...")
	}

	// 取消上下文
	mb.cancel()

	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 优雅关闭服务器
	err := mb.server.Shutdown(ctx)
	if err != nil && mb.config.EnableLog {
		log.Printf("服务器关闭失败: %v", err)
	} else if mb.config.EnableLog {
		log.Println("MiloraBot服务已停止")
	}

	return err
}

// GetBot 获取机器人实例（用于高级操作）
func (mb *MiloraBot) GetBot() *bot.Bot {
	return mb.bot
}

// GetPluginCount 获取已注册插件数量
func (mb *MiloraBot) GetPluginCount() int {
	return len(mplugin.GetPlugins())
}

// ListPlugins 列出所有已注册的插件
func (mb *MiloraBot) ListPlugins() []string {
	plugins := mplugin.GetPlugins()
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	return names
}

// DefaultConfig 返回默认配置
func DefaultConfig() *MiloraBotConfig {
	return &MiloraBotConfig{
		Port:              ":8080",
		Host:              "",
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		BotID:             123456789, // 请用户修改为自己的机器人QQ号
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableLog:         true,
		LogLevel:          "info",
		PluginDir:         "./plugins",
		AutoLoadPlugins:   true,
		PluginFilePattern: "*.so",
	}
}
