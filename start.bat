@echo off
chcp 65001 >nul
echo ========================================
echo   Message Push 系统启动 (修复版)
echo ========================================
echo.

REM 检查可执行文件
if not exist "message-push.exe" (
    echo ❌ 找不到 message-push.exe 文件！
    echo.
    echo 🔧 正在尝试编译...
    go build -o message-push.exe ./cmd/server
    if %ERRORLEVEL% NEQ 0 (
        echo ❌ 编译失败！请检查Go环境
        pause
        exit /b 1
    )
    echo ✅ 编译成功！
)

echo 🔧 设置环境变量...
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=123456
set DB_NAME=message_push
set SERVER_PORT=8080

echo.
echo 📋 当前版本包含以下SMTP修复：
echo   ✅ MIME头信息支持
echo   ✅ STARTTLS改进 (端口587)
echo   ✅ SSL/TLS优化 (端口465)
echo   ✅ 错误处理增强
echo.

echo 🚀 启动应用程序...
echo 应用将在 http://localhost:8080 运行
echo.
echo 💡 SMTP配置建议：
echo   • QQ邮箱: smtp.qq.com:587 + 授权码
echo   • Gmail: smtp.gmail.com:587 + 应用密码
echo   • 163邮箱: smtp.163.com:587 + 授权码
echo.
echo ⚠️  请不要关闭此窗口！
echo.

message-push.exe

echo.
echo 应用程序已停止
pause