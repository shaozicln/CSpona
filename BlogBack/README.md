## 文件配置：
#### config:配置
#### model:数据库存储和读写
#### api:api v1/v2/v3
#### middleware：中间件，解决跨域、jwt登陆验证
#### routes：路由接口
#### utils：公共工具包
#### upload：上传下载的目录、托管静态资源
├─  .gitignore
│  go.mod // 项目依赖
│  go.sum
│  LICENSE
│  main.go //主程序
│  README.md
│  tree.txt
│          
├─api         
├─config // 项目配置入口   
├─database  // 数据库备份文件（初始化）
├─log  // 项目日志
├─middleware  // 中间件
├─model // 数据模型层
├─routes
│      router.go // 路由入口    
├─static // 打包静态文件
│  ├─admin  // 后台管理页面 (已废弃，打包静态文件在web/admin/dist下)         
│  └─front  // 前端展示页面 (已废弃，打包静态文件在web/front/dist下)
├─upload   
├─utils // 项目公用工具库
│  │  setting.go
│  ├─errmsg   
│  └─validator         
└─web // 前端开发源码（VUECLI项目源文件)
├─admin             
└─front