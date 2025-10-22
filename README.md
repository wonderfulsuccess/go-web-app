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
