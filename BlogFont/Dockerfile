# 使用官方Node镜像，版本为16.3.0
FROM node:16.3.0 AS build

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json到工作目录
COPY package*.json ./

# 安装项目依赖
RUN npm install

# 复制项目文件到工作目录
COPY . .

# 构建项目
RUN npm run build

# 使用nginx作为服务器
FROM nginx:latest

# 复制构建后的文件到nginx的静态文件目录
COPY --from=build /app/dist /usr/share/nginx/html

# 暴露80端口
EXPOSE 80

# 运行nginx
CMD ["nginx", "-g", "daemon off;"]