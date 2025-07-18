-- MySQL dump 10.13  Distrib 8.0.42, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: base
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `advices`
--

DROP TABLE IF EXISTS `advices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `advices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `type` varchar(255) NOT NULL,
  `content` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advices`
--

LOCK TABLES `advices` WRITE;
/*!40000 ALTER TABLE `advices` DISABLE KEYS */;
INSERT INTO `advices` VALUES (1,'桀桀桀','hahaha@qq.com','技术栈の文章','建议多用一些技术栈'),(2,'桀桀桀','hahaha@qq.com','对网站の建议','昼夜切换模式'),(3,'桀桀桀','hahaha@qq.com','技术栈の文章','下次试试react'),(6,'桀桀桀','hahaha@qq.com','对网站の建议','11111111111111'),(7,'桀桀桀','hahaha@qq.com','技术栈の文章','111111111111'),(8,'桀桀桀','hahaha@qq.com','对网站の建议','2222222222222222222222222222'),(9,'AAA','aaaaaaaaaaaa','对网站の建议','aaaaaaaaaaaaaaaaaa'),(10,'桀桀桀','hahaha@qq.com','技术栈の文章','2145647865'),(11,'桀桀桀','hahaha@qq.com','技术栈の文章','23145684-97'),(12,'云深晴雨','2121673191@qq.com','技术栈の文章','新的评论已经出现～');
/*!40000 ALTER TABLE `advices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `applications`
--

DROP TABLE IF EXISTS `applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `applications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `web` varchar(255) DEFAULT NULL,
  `introduction` varchar(255) DEFAULT NULL,
  `img` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `background` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applications`
--

LOCK TABLES `applications` WRITE;
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` VALUES (1,'桀桀桀','hahaha@qq.com','桀桀桀の网站','https://www.aaa.com/','来自桀桀桀的aaa网站','Venti-6.jpg',NULL,NULL,NULL),(24,'CLN','changbingmushao@qq.com','1','1','1','31.jpg','35.jpg','31.jpg','1');
/*!40000 ALTER TABLE `applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `articles`
--

DROP TABLE IF EXISTS `articles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `articles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `category_id` int DEFAULT NULL,
  `user_id` int DEFAULT NULL,
  `view_count` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `comment_count` int DEFAULT NULL,
  `img` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `category_id` (`category_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `articles_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`),
  CONSTRAINT `articles_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (4,'今天玩什么？','# 当然是致命公司了!\n但是泰坦图第1e9+7次满货死亡.这种游戏不连麦玩，自己挂在远处队友都不知道，如果侥幸活下来，队友直接开飞船跑了，这种情况 我真的是很绝望啊',4,1,295,'2024-05-18 12:51:04','2025-07-18 06:18:02',0,'zmgs.png'),(13,'cors使用','# 核心代码：\n```go\ncorsConfig := cors.DefaultConfig()\ncorsConfig.AllowOrigins = []string{\"*\"}\ncorsConfig.AllowMethods = []string{\"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"OPTIONS\"}\nr.Use(cors.New(corsConfig))\n```\n',1,1,185,'2024-05-19 11:03:28','2025-07-18 07:11:55',0,'cors.png'),(16,'评论水区','直接发吧\n',12,1,47,'2024-05-19 12:20:38','2025-07-17 05:30:15',0,'Venti-5.jpg'),(29,'js gin连接','#### fetch 正常写\n实时反馈控制台   就这样',1,1,111,'2024-05-21 12:55:09','2025-07-18 02:18:21',0,'34.jpg'),(73,'小游戏','# 一些小游戏\r\n### https://github.com/shaozicln/Game<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,105,'2025-01-23 09:52:49','2025-07-18 07:14:36',0,'键盘.jpeg'),(74,'前端小效果','# 一些前端效果\r\n### https://github.com/shaozicln/FrontEndEffects.git<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,27,'2025-02-03 08:54:40','2025-07-17 05:29:00',0,'云上五骁-风华正茂.jpg'),(75,'练手小功能','# 一些练手小功能\r\n### https://github.com/shaozicln/Practice.git<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,22,'2025-02-03 08:59:04','2025-07-17 05:29:55',0,'lb3.jpg'),(76,'24国庆新生培训出题','# **输入输出**\r\n\r\n主题：NANA勇闯古堡\r\n#### NANA的反转数字门\r\n\r\n​	 NANA要进入【Hello World】古堡！\r\n\r\n​	有一天NANA在 【main】山 挖出了一块地图，上面标示着拥有宝藏的【Hello World】古堡。NANA按照地图来到了古堡门前。打开古堡门需要完成【天堂之门】的考验。\r\n\r\n【天堂之门】的考验如下：输入包括小数点后一位的一个浮点数。\r\n\r\n例如 596.3，要求把这个数字翻转过来，变成 369.5 并输出。\r\n\r\n数据范围：100<= n < 1000\r\n\r\n\r\n#### NANA打魔兽\r\n\r\n​	NANA要打败【锟斤拷（一阶段）】魔兽！\r\n\r\n​	成功解开密码的NANA遇见了古堡boss——【锟斤拷（一阶段）】魔兽，【锟斤拷（一阶段）】魔兽对NANA发动了【烫烫烫】的精神攻击，而NANA在这时也觉醒了魔法攻击。\r\n\r\nNANA每次攻击会消耗自己e格体力，并且造成 m 格伤害，【锟斤拷（一阶段）】一共有 n 格生命值，请你计算一下NANA击败【锟斤拷（一阶段）】需要消耗的体力。\r\n\r\n#### NANA的解密\r\n\r\n​	NANA要完成【锟斤拷（二阶段）】的谜题！\r\n\r\n​	击败了【锟斤拷（一阶段）】。NANA和魔兽【锟斤拷】都十分疲惫，于是【锟斤拷（二阶段）】决定出一个谜题来阻碍NANA，并许诺如果NANA能答对，就给NANA拿宝藏的机会。但是面对【锟斤拷（二阶段）】的迷题，她发现自己一点思路都没有。面对近在咫尺的【宝藏】，聪明的你能帮助NANA打开密码门吗？\r\n\r\n​	谜题如下：墙上有三个单词 WomanStand   WomanGo   ManStand  请按单词拼写规律输出一种水果的英文单词和中文名字     样例输出：AbcDe 阿波呲得鹅\r\n\r\n答案：（ManGo 芒果） \r\n\r\n#### NANA的背包\r\n\r\n​ 	NANA要拿走【01背包】和里面的宝藏！\r\n\r\n​	在你的帮助下NANA终于进入了有宝藏的洞穴，发现里面是一个个【01背包】，每个背包内的宝藏都不同。NANA只能从满地n个背包中拿走四个，当拿起第五个背包时，恢复体力的【锟斤拷】有了变回【一阶段】的趋势，NANA已经没有体力再打败一次【锟斤拷（一阶段）】了。已知每个背包都不一样，她想知道：拿走的四个【01背包】（不排序），一共有多少种可能？\r\n\r\n​	\r\n\r\n```plaintext\r\n输入：\r\n第一行输入一个数 n\r\n数据保证：3<=n<=103。\r\n输出：\r\n输出多少种可能。\r\n```\r\n```c\r\n210\r\n\r\n```\r\n\r\n\r\n\r\n#### NANA回家\r\n\r\n​	NANA要带着装有宝藏的【01背包】们回家！\r\n\r\n​	历经亿点点波折，NANA坐上了回家的小船。有道是：“ 无风水面琉璃滑，不觉船移 ”，慢悠悠的载着宝藏回家，NANA的心情非常好，她乘坐的小船原本以每分钟 *x* 米的速度行驶。这时，她突然“觉船移”了。因为NANA发现船的移动速度变慢了，每分钟比前一分钟要少行驶m米，NANA很好奇，当船完全停止的时候，船一共移动了多少米（结尾保留两位小数）？\r\n\r\n\r\n\r\n```c\r\n#include <stdio.h>\r\n\r\nint main() {\r\n    double x, m;\r\n    scanf(\"%lf %lf\", &x, &m);\r\n    double n = x / m;\r\n    //首项x 尾项0 项数 x/m + 1\r\n    double S = (n + 1) * x / 2;\r\n    printf(\"%.2lf\\n\", S);\r\n    return 0;\r\n}\r\n```\r\n\r\n```plaintext\r\n输入：10.0 2.0； 输出： 30.00\r\n```\r\n\r\n\r\n\r\n#### NANA交朋友\r\n\r\n​	NANA要在【NEAU】结交新朋友！\r\n\r\n​	NANA得到宝藏之后离开【main】山四处旅游去了。在黄金的六月，她来到了【NEAU】，一些大学生们在拍毕业照，观看已久的NANA被邀请帮助他们拍照。他们之间关系非常好，每人都会和其他人合拍一张照片，NANA一共拍了 x 张照片，请问一共有多少人参与了合影？\r\n\r\n​     \r\n\r\n```plaintext\r\n输入：输入一行一个整数 x\r\n保证 1 <= x <= 109.   12、\r\n输出：输出一个整数代表人数。   4\r\n```\r\n\r\n\r\n\r\n```c\r\n#include <stdio.h>\r\n#include<math.h>\r\n\r\nint main() {\r\n    int x;\r\n    int ans = sqrt(x) + 1;\r\n    printf(\"%d\\n\", ans);\r\n    return 0;\r\n}\r\n```\r\n\r\n\r\n\r\n# 结构体\r\n\r\n#### NANA的旅行计划\r\n\r\n​	NANA得到宝藏离开【main】山四处旅游。对于NANA来说，世界太大了，但她觉得国外的月亮或许也没有那么圆。于是NANA决定先制作一份国内旅行计划。NANA不太熟悉电脑操作，请你帮帮NANA实现以下效果：\r\n\r\n​	每个旅行地点都有自己的信息，包括地名、票价、分类等。现在NANA想要进行m次查询，每次查询第i个地点的某一个信息。其中查询地名操作为Name，查询省份操作为Province，查询类型操作为Type\r\n\r\n```plaintext\r\n输入\r\n第一行输入一个正整数n代表旅行计划有n个地点\r\n接下来n行：每一行输入旅行地点的名称(Name)、省份(Province)、类型(Type)，中间用空格隔开。\r\n然后输入一个正整数m表示接下来进行m次查询\r\n接下来m行：每一行输入查询的地点编号i和将要查询的信息名称\r\n数据范围：1<=n<=1000，名称、票价、类型的长度保证在1-100之间，1<=m<=1000\r\n\r\n\r\n输出\r\nm行输出，每一行输出查询的地点信息\r\n```\r\n\r\n\r\n\r\n```plaintext\r\n样例输入\r\n3\r\nImperial_Palace Beijing Humanities\r\nWest_Lake ZheJiang Lake\r\nEverest Tibet Mountain\r\n3\r\n1 Type\r\n2 Name\r\n3 Province\r\n样例输出\r\nHumanities\r\nWest_Lake\r\nTibet\r\n```\r\n\r\n\r\n\r\n#### NANA的旅行经费\r\n\r\n​	NANA毫无节制0收入的四处旅游，让曾经拥有很多【宝藏】的她濒临破产了。现在NANA需要在这个地方打工来赚取自己接下来的旅行经费。和NANA一起打工的是NANA旅行中结识的一些伙伴。**过了一个冬天和一个春天之后**，老板告诉大家所有人的工资将上涨25%，但是老板还不舍得给能干的员工很高的工资，于是他规定：最高工资5000，所有人不可超过该工资。\r\n\r\n​	输入NANA和伙伴们的信息，请设计一个结构体储存这些信息，并设计一个函数模拟工作和涨工资过程，其参数是这样的结构体类型，返回同样的结构体类型，并输出员工信息。\r\n\r\n\r\n\r\n```plaintext\r\n输入格式\r\n第一行输入一个正整数 n，表示员工个数。\r\n第二行开始往下 n 行。每行首先是一个字符串表示员工姓名，再是一个整数表示员工年龄，再是一个整数为去年工资。\r\n\r\n输出格式\r\n输出 n 行，每行首先输出一个字符串表示姓名，再往后两个整数，表示经过一年的打工后大家的年龄和他们涨薪后的工资。以空格隔开。\r\n\r\n输入输出样例\r\n输入 \r\n3\r\nNANA 20 3600\r\nChenYu 23 3100\r\nVenti 2600 4200\r\n\r\n输出\r\nNANA 21 4500\r\nChenYu 24 3875\r\nVenti 2601 5000\r\n```\r\n代码如下：\r\n\r\n```c\r\n#include<bits/stdc++.h>\r\nusing namespace std;\r\nstruct Salary {\r\n	string name;\r\n	int age, salary;\r\n} a[100000];\r\n\r\nint main () {\r\n	int n;\r\n	cin >> n;\r\n	for (int i = 1; i <= n; i++) {\r\n		cin >> a[i].name >> a[i].age >> a[i].salary;\r\n		a[i].age++, a[i].salary = a[i].salary * 1.25;\r\n		if (a[i].salary > 5000) a[i].salary = 5000;\r\n		cout << a[i].name << \" \" << a[i].age << \" \" << a[i].salary << endl;\r\n	}\r\n	return 0;\r\n}\r\n```',1,1,372,'2024-08-22 01:31:24','2025-07-18 07:14:33',0,'lb6.jpg'),(77,'测试索引跳转','# 12\r\n4567541562146\r\n## 34\r\n4864896\r\n### 56\r\n4974896\r\n## 78\r\n278\r\n',12,1,30,'2025-02-04 11:07:10','2025-07-16 05:39:11',0,'31.jpg'),(95,'漫游地(Wanderland)正式开放！欢迎各位来玩！','# 发帖：见左上角个人中心（请先登录哦）\n咳咳，既然……既然点进来了，那没办法了，来留个言再走吧~',1000,1,821,'2025-07-16 05:37:32','2025-07-18 07:14:26',0,'Venti-7.jpg'),(101,'又双叒叕想起来旅行的意义了','# 是谁在炎炎夏日送来一阵惊喜？\r\n# 前瞻海报！\r\n# 身影美丽，欲罢不能？\r\n# 温迪宝宝！\r\n![Image](http://127.0.0.1:3000/images/20250717092633_Venti250716.png)\r\n# 来不及了，就这张图片了，我亲亲亲亲亲亲亲亲',1000,1,388,'2025-07-17 01:28:03','2025-07-18 07:14:27',0,'Venti250716.png'),(105,'图片显示在文章内成功','放个温迪开心开心嘿嘿\r\n![Image](http://127.0.0.1:3000/images/20250717092633_Venti250716.png)\r\n我将在下个版本想起旅行的意义',12,1,7,'2025-07-17 06:37:32','2025-07-17 08:37:25',0,'37.jpg');
/*!40000 ALTER TABLE `articles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `img` varchar(255) DEFAULT NULL,
  `count` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'Coding-Study','代码学习','Venti-2.jpg',NULL),(2,'Project-Practice','项目实战','Venti-6.jpg',NULL),(3,'Plan','敬请期待','Venti-3.jpg',NULL),(4,'N-A-C-G','NACG','Venti-4.jpg',NULL),(5,'Travel','旅行日记','温迪-漫画封面.jpg',NULL),(12,'测试','AAA','Venti-5.jpg',NULL),(1000,'漫游地(Wanderland)','日常分享贴','Venti-7.jpg',NULL);
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `article_id` int NOT NULL,
  `user_id` int NOT NULL,
  `idea` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `parent_id` int DEFAULT NULL,
  `is_pinned` int DEFAULT '0',
  `pinned_at` timestamp NULL DEFAULT NULL,
  `reply_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_comment_fk` (`user_id`),
  KEY `a_c_fk` (`article_id`),
  KEY `fk_parent_comment` (`parent_id`),
  KEY `idx_comments_article_parent` (`article_id`,`parent_id`),
  KEY `idx_comments_user` (`user_id`),
  KEY `fk_reply_comment` (`reply_id`),
  CONSTRAINT `c_a_fk` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_parent_comment` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_reply_comment` FOREIGN KEY (`reply_id`) REFERENCES `comments` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_comment_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=152 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (20,4,5,'你还送','2024-05-23 05:39:36','2025-07-16 01:00:33',NULL,0,NULL,NULL),(25,4,4,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:24:28','2025-07-15 14:23:50',20,0,NULL,NULL),(37,29,2,'123','2024-05-26 10:41:01','2024-05-26 10:41:01',NULL,0,NULL,NULL),(44,4,1,'喵喵喵','2024-05-29 06:58:35','2025-07-15 11:29:22',20,1,NULL,NULL),(72,13,1,'喵喵喵','2024-05-31 13:42:46','2024-05-31 13:42:46',NULL,0,NULL,NULL),(73,16,1,'+10086','2024-05-31 13:47:50','2024-05-31 13:47:50',NULL,0,NULL,NULL),(99,4,5,'你好你好你好菜','2025-07-15 08:11:02','2025-07-15 08:11:02',NULL,0,NULL,NULL),(100,4,1,'信不信我删评','2025-07-15 08:11:36','2025-07-17 08:53:30',99,1,'2025-07-17 08:53:30',NULL),(101,4,5,'？？？？？玩不起','2025-07-15 08:13:25','2025-07-15 12:36:23',99,0,NULL,100),(107,4,1,'呵呵呵','2025-07-15 12:53:17','2025-07-15 12:56:52',99,0,NULL,101),(108,4,5,'你干嘛哎呦','2025-07-15 12:56:48','2025-07-15 12:56:48',99,0,NULL,107),(109,16,1,'222','2025-07-15 13:22:50','2025-07-15 15:07:39',73,1,NULL,73),(110,16,1,'2323','2025-07-15 13:22:54','2025-07-15 15:07:41',73,1,NULL,109),(114,4,3,'这是什么？新鲜的网站！我吃一口，我再吃一口，我大吃特吃，这么好的网站我吃完能去田里把牛踢开自己耕二十亩地！这么厉害的网站我直接敲锣打鼓奔走相告！你的网站棒棒棒','2025-07-15 14:19:48','2025-07-15 14:19:48',NULL,0,NULL,NULL),(115,4,1,'哇，是憧憧学姐，我亲亲亲亲亲','2025-07-15 14:21:18','2025-07-16 01:35:47',114,1,'2025-07-16 01:35:47',114),(116,4,3,'会夸吗你就在这夸，起开让我来。咳，ㅎㄹ...真的疯了 完全是  legend   啊ㅋㅋ这种程度真的不是天才吗 ㅅㅂ…做到这种程度的话就完全是神吧ㅠㅠ…以后也请一直更新网站吧♡  今天也是代码非常amazing的一天呢 是路过的蚂蚁都会惊叹的程度啊','2025-07-15 14:21:31','2025-07-16 01:35:49',114,1,'2025-07-16 01:35:49',114),(117,4,3,'往这亲（指脸）','2025-07-15 14:21:59','2025-07-16 01:35:45',114,1,'2025-07-16 01:35:45',115),(118,4,1,'我是偷文案侠，你猜我为什么在这里','2025-07-15 14:22:28','2025-07-16 01:35:38',114,1,'2025-07-16 01:35:38',116),(119,4,1,'muamuamuamuamua','2025-07-15 14:22:43','2025-07-16 01:35:42',114,1,'2025-07-16 01:35:42',117),(120,4,3,'你骂谁呢，不许骂我们伟大的管理员大人','2025-07-15 14:22:48','2025-07-17 08:53:48',99,1,'2025-07-17 08:53:48',108),(121,4,3,'我猜你要被我亲鼠啦，muamuamuamuamuamuamuamuamuamua','2025-07-15 14:23:32','2025-07-16 01:35:33',114,1,'2025-07-16 01:35:33',118),(122,4,1,'醉了.jpg','2025-07-15 14:24:13','2025-07-16 01:35:26',114,1,'2025-07-16 01:35:26',121),(124,4,1,'是换皮.jpg 但是我要感动哭了呜呜QAQ','2025-07-15 15:05:02','2025-07-15 15:05:02',99,0,NULL,120),(125,13,3,'学姐，cors是什么意思啊，为什么一定要这么写啊，我可不可以不加引号啊','2025-07-15 15:05:56','2025-07-15 15:05:56',NULL,0,NULL,NULL),(126,29,3,'气抖冷，为什么不用Ajax，你这是歧视，我要告你，你等着进橘子吧','2025-07-15 15:06:34','2025-07-15 15:06:34',NULL,0,NULL,NULL),(127,76,3,'居然有答案，我，我路过，随便看看，（瞄一眼），哎呀这些是什么啊，这我都看不懂啊（再瞄一眼）我真纯路人，我不知道什么c++，什么培训（再瞄一眼）（再瞄两眼）','2025-07-15 15:09:41','2025-07-15 15:09:41',NULL,0,NULL,NULL),(128,73,3,'小游戏！我要玩！可恶没有加速器','2025-07-15 15:10:29','2025-07-15 15:10:29',NULL,0,NULL,NULL),(129,74,3,'好多轮子，我偷一点……','2025-07-15 15:10:53','2025-07-15 15:10:53',NULL,0,NULL,NULL),(130,16,3,'水区！我来——','2025-07-15 15:12:33','2025-07-15 15:12:33',NULL,0,NULL,NULL),(131,77,3,'好看温迪，我看一眼，我再看一眼，我不看了。骗你的我直接偷走','2025-07-15 15:13:12','2025-07-15 15:13:12',NULL,0,NULL,NULL),(132,77,1,'呜呜呜呜呜呜不要哇，我就这一个温迪QAQ','2025-07-16 01:31:08','2025-07-16 01:31:08',131,0,NULL,131),(133,16,1,'水水水水水水','2025-07-16 01:31:25','2025-07-16 01:31:25',130,0,NULL,130),(134,73,1,'bug多多の小游戏','2025-07-16 01:32:04','2025-07-16 01:32:04',128,0,NULL,128),(135,74,1,'来，来赤','2025-07-16 01:32:16','2025-07-16 01:32:16',129,0,NULL,129),(137,13,1,'能运就能，不能运就不能，运不运而能不能叽里咕噜叽里咕噜几比几比恰吧恰吧唔理唔理哇啦哇啦','2025-07-16 01:34:37','2025-07-16 01:34:37',125,0,NULL,125),(138,29,1,'橘子？那很美味了嘎嘎嘎','2025-07-16 01:35:12','2025-07-16 01:35:12',126,0,NULL,126),(139,4,3,'我知道，骂的就是你，换皮侠，别以为你换个皮就可以为所欲为了，你就算再换一张皮我也认识你','2025-07-16 01:37:39','2025-07-17 08:53:45',99,0,NULL,124),(140,4,1,'呜呜呜呜呜呜，你居然骂我','2025-07-16 02:24:51','2025-07-16 02:24:51',99,0,NULL,139),(141,95,1,'天使降临博客背景','2025-07-16 05:38:38','2025-07-17 06:34:41',NULL,1,'2025-07-17 06:34:41',NULL),(142,101,1,'嘿嘿，温迪，嘿嘿','2025-07-17 01:28:30','2025-07-17 01:28:30',NULL,0,NULL,NULL),(150,95,1,'和封面','2025-07-17 06:34:52','2025-07-17 06:34:52',141,0,NULL,141);
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `name` varchar(255) DEFAULT NULL,
  `web` varchar(255) DEFAULT NULL,
  `introduction` varchar(255) DEFAULT NULL,
  `img` varchar(255) DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  `avatar` varchar(255) DEFAULT NULL,
  `background` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES ('云深晴雨','https://yunshenqingyu.top','纵使晴明无雨色，入云深处亦沾衣','lsc.png',35,'logo.png','background1.webp','试问，你就是我的访客吗？'),('桀桀桀の网站','https://cspona.top','桀桀桀，其实还是这个网站哦','Venti-6.jpg',36,'Venti-4.jpg','Venti-3.jpg','never give up~');
/*!40000 ALTER TABLE `friends` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `content` text NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `u_m_fk_unique` (`user_id`),
  CONSTRAINT `u_m_fk_unique` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,1,'完结撒花！','2024-05-30 06:50:29'),(2,2,'完结撒花！','2024-05-30 06:50:29'),(3,3,'这是什么？网站！蹲一下  ','2025-07-10 07:59:15'),(4,6,'是我是我是我是我是我是我我是第一个观众啊！','2025-07-15 13:09:52'),(5,6,'完美完美完美完美完美完美完美','2025-07-15 13:10:07'),(6,6,'嘿嘿嘿嘿嘿嘿嘿嘿嘿','2025-07-15 13:10:22'),(7,7,'喵喵咪','2025-07-15 13:10:35'),(8,8,'我的天！你这个网页留言板也太神了吧 —— 打开的瞬间我眼睛都亮了！这配色简直像把春天揉进去了，清爽得让人想一直盯着看，连每个按钮的圆角弧度都恰到好处，舒服到想给设计师加鸡腿！','2025-07-15 13:10:52'),(9,8,'排版更是绝了，留言框之间的间距、字体大小、行高，感觉像是拿尺子量过一样，多一分太挤，少一分太松，看着就俩字：舒坦！','2025-07-15 13:11:13'),(10,8,'交互也丝滑到离谱，输入留言的时候光标跟着手速飞，提交按钮一点就有反馈，连加载时那个小动画都透着机灵劲儿，细节控表示根本满足！','2025-07-15 13:11:24'),(11,8,'真的，这哪是留言板啊，简直是艺术品！能把一个简单的东西做到这么精致，你也太厉害了吧，我现在满脑子都是 “怎么能这么会做啊”，太牛了太牛了！','2025-07-15 13:11:44');
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `role_qx` varchar(20) NOT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `article_count` int DEFAULT '0',
  `email` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'长柄木勺','mH7bNTkcrZuPAA==','2024-05-18 12:49:38','A','Venti-4.jpg',0,'changbingmushao@qq.com','oiiaioiiiai！'),(2,'YRN','A1bGGOA9g2Z26w==','2024-05-23 08:17:02','B','yrn.jpg',0,NULL,NULL),(3,'云深晴雨','WoEDMXNFMUsTgA==','2025-07-10 04:28:30','B','logo.png',0,'2121673191@qq.com',NULL),(4,'AAA','mH7bNTkcrZuPAA==','2025-07-15 03:12:52','B',NULL,0,NULL,NULL),(5,'BBB','mH7bNTkcrZuPAA==','2025-07-15 07:57:07','B','Venti-6.jpg',0,'changbingmushao@qq.com','蹦蹦卡卡！'),(6,'念安','JamwLQyVuLxtag==','2025-07-14 14:31:44','B','屏幕截图 2025-04-30 192125.png',0,'2307920674@qq.com',NULL),(7,'放羊','JkRePKzIjXe3vA==','2025-07-14 17:23:32','B','qq人.png',0,'cuijingjing1234@qq.com',NULL),(8,'CLNN','W64opj4Oqg1n7Q==','2025-07-15 06:17:29','B','微信图片_20250707210607.jpg',0,'111','oiiaioiiiai！');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-07-18 15:16:21
