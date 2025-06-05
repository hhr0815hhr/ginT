## GIN Web服务模板

核心功能模块实现：
   - 数据库访问层(model/mysql)
   - 业务逻辑层(logic)
   - 控制器层(controller)
   - 中间件(middleware)
   - 工具函数(util)
   - 定时任务(cron)
   - 队列系统(queue)
   - 国际化(i18n)

3. 基础设施集成：
   - Redis缓存
   - MySQL数据库
   - 邮件服务
   - Google登录
   - Airwallex支付

4. 命令行工具：
   - 服务启动
   - 定时任务
   - 队列消费
   - 代码生成

5. 辅助功能：
   - 日志系统
   - 配置管理
   - 响应处理
   - 加密工具
   - 随机字符串生成

项目采用模块化设计，使用wire实现依赖注入，遵循清晰的代码分层架构。各模块职责明确，便于后续扩展和维护。

主要技术栈：
- Web框架: Gin
- ORM: GORM
- 依赖注入: Wire
- 配置管理: Viper
- 日志: Logrus
- 定时任务: Cron
- 队列: Redis/Memory