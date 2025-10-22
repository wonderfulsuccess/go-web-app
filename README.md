后端 go 前端 react-ts
开发桌面 web 应用程序模板
back 为后端目录
front 为前端目录

## 项目介绍

形成一套利用 go 作为后端，react-ts 作为前端，开发桌面应用程序的模板

## 功能特性

1. 使用 gin 作为后端 webserver，webserver 相关的代码路由、API、都定义在 back/webserver 路径下。静态资源存放在 back/webserver/dist 路径下。dist 由前端 front 项目生成
2. webserver 支持 websocket，提供一个专门的函数，用于发送和接受 websocket 消息。发送函数发送的数据包含，发送方、接受方、时间戳、消息类型，消息内容。消息内容为 json 格式。前端有一个专门的 API 用于订阅 websocket 消息，是否使用改消息由具体组件根据接受方和发送方综合判断。
3. 后端使用 gorm 操作数据库，支持 mysql、postgres、sqlite，默认使用 sqlite，支持一行代码切换数据库。数据表定义在 back/model 目录下，一个数据表一个 go 代码文件。每个数据表对应的增删改查等操作定义在 back/controller 目录下，一个数据表对应一个 go 代码文件。
4. 前端支持 tailwindcss、shadcn-ui、react-icons。创建一个建议的后台管理模板的单页应用 SPA，顶部是菜单，可以切换不同的页面。每个菜单有自己的 url 路径。页面支持亮暗主题，可以通过按钮主动切换，也可以跟随系统自动切换。

## 快速开始

### 后端（Go）

1. 进入目录并安装依赖
   ```bash
   cd back
   go mod tidy
   ```
2. 运行开发服务器
   ```bash
   go run ./app
   ```
3. 相关配置（可选）通过环境变量控制：
   - `SERVER_PORT`：HTTP 服务端口，默认 `8080`
   - `STATIC_DIR`：静态资源目录，默认 `back/webserver/dist`
   - `DB_TYPE`：数据库类型，可选 `sqlite`（默认）/`mysql`/`postgres`
   - `DB_DSN`：数据库连接串。使用 `sqlite` 时默认生成 `back/data/app.db`
   - `DB_LOG_SQL`：是否输出 Gorm SQL 日志，默认关闭，设置为 `true` 启用

> 数据模型存放在 `back/model`，控制器在 `back/controller`。`back/webserver` 统一注册路由、API 与 WebSocket 入口，同时负责分发 `front` 构建出的静态资源。

### 前端（React + Vite）

1. 安装依赖
   ```bash
   cd front
   npm install
   ```
2. 启动开发服务
   ```bash
   npm run dev
   ```
3. 构建产物（会输出到 `back/webserver/dist`，供 Gin 静态服务使用）
   ```bash
   npm run build
   ```

   完成构建后，直接启动后端服务即可通过 [http://localhost:8080/](http://localhost:8080/) 访问页面。

> TailwindCSS、shadcn/ui 与 react-icons 已预配置，可直接在 `src` 下按需引入。前端路由基于 `react-router-dom`，默认包含仪表盘、用户管理、系统设置三个页面，并支持亮暗主题切换。

### WebSocket 使用

- 后端：通过 `webserver.Hub` 的 `SendMessage` 方法发送标准化消息（包含发送方、接收方、时间戳、消息类型、JSON 消息体）。所有来自客户端的消息也会统一进入 `Hub.Incoming()` 便于二次处理。
- 前端：`src/api/websocket.ts` 提供 `connectWebSocket` 与 `subscribeToMessages` 方法集中管理连接与订阅，组件只需调用订阅函数即可接收实时推送，同时可以使用 `sendMessage` 在需要时主动发送消息。

## 调试建议

- 后端：`go test ./...` 可快速校验代码是否可以正常编译。
- 前端：`npm run build` 会在构建阶段执行 TypeScript 类型检查；如需格式或质量检查可扩展 `npm run lint`。

> 当前 Vite 依赖需要 Node.js >= 20.19，若本地低于该版本会在构建阶段提示升级。
