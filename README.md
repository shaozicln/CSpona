数据库：库名和密码在BlogBack/config/config.ini文件
注：
- mysql建用户命令行：
``` CREATE USER 'svb'@'host' IDENTIFIED BY 'svb2861';```
- 授权命令行：
```GRANT ALL PRIVILEGES ON base.* TO'svb'@'host';```
- 刷新权限：
```FLUSH PRIVILEGES;```
- 克隆、导入数据库：
```CREATE DATABASE base;```
```mysql -u svb svb2861 base < base.sql;```
- 检验数据库：
```use base;```
```show tables;```
```select * from <table_name>;```

vue前端网页:
```npm install```
```npm run dev```
```npm run build```