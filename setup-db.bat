@echo off
chcp 65001 >nul
echo ========================================
echo   配置本地 PostgreSQL 数据库
echo ========================================
echo.

REM 设置 PostgreSQL 路径
set PSQL_PATH="C:\Program Files\PostgreSQL\17\bin\psql.exe"

echo 检测到 PostgreSQL 17 已安装
echo 路径: C:\Program Files\PostgreSQL\17\bin\
echo.

echo 请确保您的本地 PostgreSQL 服务正在运行
echo.
echo 默认连接信息：
echo   主机: localhost
echo   端口: 5432
echo   用户: postgres
echo.

set /p PGPASSWORD=请输入 PostgreSQL 密码: 

echo.
echo [1/4] 测试数据库连接...
%PSQL_PATH% -U postgres -h localhost -p 5432 -c "SELECT version();" -t >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 数据库连接成功
) else (
    echo ❌ 数据库连接失败，请检查密码
    pause
    exit /b 1
)

echo.
echo [2/4] 创建数据库...
%PSQL_PATH% -U postgres -h localhost -p 5432 -c "CREATE DATABASE message_push;" 2>nul
if %errorlevel% equ 0 (
    echo ✅ 数据库创建成功
) else (
    echo ⚠️  数据库可能已存在，继续...
)

echo.
echo [3/4] 执行数据库迁移...
%PSQL_PATH% -U postgres -h localhost -p 5432 -d message_push -f migrations/init.sql
if %errorlevel% equ 0 (
    echo ✅ 数据表创建成功
) else (
    echo ❌ 数据表创建失败
    pause
    exit /b 1
)

echo.
echo [4/4] 验证数据库...
echo 数据库表列表：
%PSQL_PATH% -U postgres -h localhost -p 5432 -d message_push -c "\dt"

echo.
echo ========================================
echo   数据库配置完成！
echo ========================================
echo.
echo 数据库信息：
echo   数据库名: message_push
echo   主机: localhost
echo   端口: 5432
echo   用户: postgres
echo.
echo 创建的表：
echo   - groups (群组信息)
echo   - members (成员信息)
echo   - messages (消息记录)
echo   - push_logs (推送日志)
echo.
pause