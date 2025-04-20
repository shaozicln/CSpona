-- MySQL dump 10.13  Distrib 8.0.37, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: base
-- ------------------------------------------------------
-- Server version	8.0.37

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
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advices`
--

LOCK TABLES `advices` WRITE;
/*!40000 ALTER TABLE `advices` DISABLE KEYS */;
INSERT INTO `advices` VALUES (1,'桀桀桀','hahaha@qq.com','技术栈の文章','建议多用一些技术栈'),(2,'桀桀桀','hahaha@qq.com','对网站の建议','昼夜切换模式'),(3,'桀桀桀','hahaha@qq.com','技术栈の文章','下次试试react'),(6,'桀桀桀','hahaha@qq.com','对网站の建议','11111111111111'),(7,'桀桀桀','hahaha@qq.com','技术栈の文章','111111111111'),(8,'桀桀桀','hahaha@qq.com','对网站の建议','2222222222222222222222222222'),(9,'AAA','aaaaaaaaaaaa','对网站の建议','aaaaaaaaaaaaaaaaaa'),(10,'桀桀桀','hahaha@qq.com','技术栈の文章','2145647865'),(11,'桀桀桀','hahaha@qq.com','技术栈の文章','23145684-97');
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applications`
--

LOCK TABLES `applications` WRITE;
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` VALUES (1,'桀桀桀','hahaha@qq.com','桀桀桀の网站','https://www.aaa.com/','来自桀桀桀的aaa网站','Venti-6.jpg'),(3,'桀桀桀','hahaha@qq.com','111111111','111111111111','11111111111111111','35.jpg'),(5,'CLN','changbingmushao@qq.com','qqqqqq','qqqqqqqqqq','qqq','35.jpg'),(7,'CLN','changbingmushao@qq.com','测试图片用网站申请','ww.test.com','123456798','星星5.jpg'),(9,'CLN','changbingmushao@qq.com','1234','123445','1431241542','屏幕截图 2024-06-01 224328.png'),(10,'JWW','','1245','124','1243分得的','屏幕截图 2024-06-24 104721.png');
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
) ENGINE=InnoDB AUTO_INCREMENT=78 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (4,'今天玩什么？','# 当然是致命公司了!\n但是泰坦图第1e9+7次满货死亡.这种游戏不连麦玩，自己挂在远处队友都不知道，如果侥幸活下来，队友直接开飞船跑了，这种情况 我真的是很绝望啊',4,1,116,'2024-05-18 12:51:04','2025-02-04 01:21:34',0,'zmgs.png'),(13,'cors使用','# 核心代码：\n```go\ncorsConfig := cors.DefaultConfig()\ncorsConfig.AllowOrigins = []string{\"*\"}\ncorsConfig.AllowMethods = []string{\"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"OPTIONS\"}\nr.Use(cors.New(corsConfig))\n```\n',1,1,75,'2024-05-19 11:03:28','2025-02-04 13:20:45',0,'cors.png'),(16,'评论水区','直接发吧\n',12,1,26,'2024-05-19 12:20:38','2025-02-03 13:22:13',0,'Venti-5.jpg'),(29,'js gin链接','#### fetch 正常写\n实时反馈控制台   就这样',1,1,17,'2024-05-21 12:55:09','2025-02-04 13:10:55',0,'34.jpg'),(73,'小游戏','# 一些小游戏\r\n### https://github.com/shaozicln/Game<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,93,'2025-01-23 09:52:49','2025-02-04 13:20:47',0,'键盘.jpeg'),(74,'前端小效果','# 一些前端效果\r\n### https://github.com/shaozicln/FrontEndEffects.git<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,22,'2025-02-03 08:54:40','2025-02-04 13:20:39',0,'云上五骁-风华正茂.jpg'),(75,'练手小功能','# 一些练手小功能\r\n### https://github.com/shaozicln/Practice.git<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,13,'2025-02-03 08:59:04','2025-02-04 13:20:28',0,'lb3.jpg'),(76,'24国庆新生培训出题','# **输入输出**\r\n\r\n主题：NANA勇闯古堡\r\n#### NANA的反转数字门\r\n\r\n​	 NANA要进入【Hello World】古堡！\r\n\r\n​	有一天NANA在 【main】山 挖出了一块地图，上面标示着拥有宝藏的【Hello World】古堡。NANA按照地图来到了古堡门前。打开古堡门需要完成【天堂之门】的考验。\r\n\r\n【天堂之门】的考验如下：输入包括小数点后一位的一个浮点数。\r\n\r\n例如 596.3，要求把这个数字翻转过来，变成 369.5 并输出。\r\n\r\n数据范围：100<= n < 1000\r\n\r\n\r\n#### NANA打魔兽\r\n\r\n​	NANA要打败【锟斤拷（一阶段）】魔兽！\r\n\r\n​	成功解开密码的NANA遇见了古堡boss——【锟斤拷（一阶段）】魔兽，【锟斤拷（一阶段）】魔兽对NANA发动了【烫烫烫】的精神攻击，而NANA在这时也觉醒了魔法攻击。\r\n\r\nNANA每次攻击会消耗自己e格体力，并且造成 m 格伤害，【锟斤拷（一阶段）】一共有 n 格生命值，请你计算一下NANA击败【锟斤拷（一阶段）】需要消耗的体力。\r\n\r\n#### NANA的解密\r\n\r\n​	NANA要完成【锟斤拷（二阶段）】的谜题！\r\n\r\n​	击败了【锟斤拷（一阶段）】。NANA和魔兽【锟斤拷】都十分疲惫，于是【锟斤拷（二阶段）】决定出一个谜题来阻碍NANA，并许诺如果NANA能答对，就给NANA拿宝藏的机会。但是面对【锟斤拷（二阶段）】的迷题，她发现自己一点思路都没有。面对近在咫尺的【宝藏】，聪明的你能帮助NANA打开密码门吗？\r\n\r\n​	谜题如下：墙上有三个单词 WomanStand   WomanGo   ManStand  请按单词拼写规律输出一种水果的英文单词和中文名字     样例输出：AbcDe 阿波呲得鹅\r\n\r\n答案：（ManGo 芒果） \r\n\r\n#### NANA的背包\r\n\r\n​ 	NANA要拿走【01背包】和里面的宝藏！\r\n\r\n​	在你的帮助下NANA终于进入了有宝藏的洞穴，发现里面是一个个【01背包】，每个背包内的宝藏都不同。NANA只能从满地n个背包中拿走四个，当拿起第五个背包时，恢复体力的【锟斤拷】有了变回【一阶段】的趋势，NANA已经没有体力再打败一次【锟斤拷（一阶段）】了。已知每个背包都不一样，她想知道：拿走的四个【01背包】（不排序），一共有多少种可能？\r\n\r\n​	\r\n\r\n```plaintext\r\n输入：\r\n第一行输入一个数 n\r\n数据保证：3<=n<=103。\r\n输出：\r\n输出多少种可能。\r\n```\r\n```c\r\n210\r\n\r\n```\r\n\r\n\r\n\r\n#### NANA回家\r\n\r\n​	NANA要带着装有宝藏的【01背包】们回家！\r\n\r\n​	历经亿点点波折，NANA坐上了回家的小船。有道是：“ 无风水面琉璃滑，不觉船移 ”，慢悠悠的载着宝藏回家，NANA的心情非常好，她乘坐的小船原本以每分钟 *x* 米的速度行驶。这时，她突然“觉船移”了。因为NANA发现船的移动速度变慢了，每分钟比前一分钟要少行驶m米，NANA很好奇，当船完全停止的时候，船一共移动了多少米（结尾保留两位小数）？\r\n\r\n\r\n\r\n```c\r\n#include <stdio.h>\r\n\r\nint main() {\r\n    double x, m;\r\n    scanf(\"%lf %lf\", &x, &m);\r\n    double n = x / m;\r\n    //首项x 尾项0 项数 x/m + 1\r\n    double S = (n + 1) * x / 2;\r\n    printf(\"%.2lf\\n\", S);\r\n    return 0;\r\n}\r\n```\r\n\r\n```plaintext\r\n输入：10.0 2.0； 输出： 30.00\r\n```\r\n\r\n\r\n\r\n#### NANA交朋友\r\n\r\n​	NANA要在【NEAU】结交新朋友！\r\n\r\n​	NANA得到宝藏之后离开【main】山四处旅游去了。在黄金的六月，她来到了【NEAU】，一些大学生们在拍毕业照，观看已久的NANA被邀请帮助他们拍照。他们之间关系非常好，每人都会和其他人合拍一张照片，NANA一共拍了 x 张照片，请问一共有多少人参与了合影？\r\n\r\n​     \r\n\r\n```plaintext\r\n输入：输入一行一个整数 x\r\n保证 1 <= x <= 109.   12、\r\n输出：输出一个整数代表人数。   4\r\n```\r\n\r\n\r\n\r\n```c\r\n#include <stdio.h>\r\n#include<math.h>\r\n\r\nint main() {\r\n    int x;\r\n    int ans = sqrt(x) + 1;\r\n    printf(\"%d\\n\", ans);\r\n    return 0;\r\n}\r\n```\r\n\r\n\r\n\r\n# 结构体\r\n\r\n#### NANA的旅行计划\r\n\r\n​	NANA得到宝藏离开【main】山四处旅游。对于NANA来说，世界太大了，但她觉得国外的月亮或许也没有那么圆。于是NANA决定先制作一份国内旅行计划。NANA不太熟悉电脑操作，请你帮帮NANA实现以下效果：\r\n\r\n​	每个旅行地点都有自己的信息，包括地名、票价、分类等。现在NANA想要进行m次查询，每次查询第i个地点的某一个信息。其中查询地名操作为Name，查询省份操作为Province，查询类型操作为Type\r\n\r\n```plaintext\r\n输入\r\n第一行输入一个正整数n代表旅行计划有n个地点\r\n接下来n行：每一行输入旅行地点的名称(Name)、省份(Province)、类型(Type)，中间用空格隔开。\r\n然后输入一个正整数m表示接下来进行m次查询\r\n接下来m行：每一行输入查询的地点编号i和将要查询的信息名称\r\n数据范围：1<=n<=1000，名称、票价、类型的长度保证在1-100之间，1<=m<=1000\r\n\r\n\r\n输出\r\nm行输出，每一行输出查询的地点信息\r\n```\r\n\r\n\r\n\r\n```plaintext\r\n样例输入\r\n3\r\nImperial_Palace Beijing Humanities\r\nWest_Lake ZheJiang Lake\r\nEverest Tibet Mountain\r\n3\r\n1 Type\r\n2 Name\r\n3 Province\r\n样例输出\r\nHumanities\r\nWest_Lake\r\nTibet\r\n```\r\n\r\n\r\n\r\n#### NANA的旅行经费\r\n\r\n​	NANA毫无节制0收入的四处旅游，让曾经拥有很多【宝藏】的她濒临破产了。现在NANA需要在这个地方打工来赚取自己接下来的旅行经费。和NANA一起打工的是NANA旅行中结识的一些伙伴。**过了一个冬天和一个春天之后**，老板告诉大家所有人的工资将上涨25%，但是老板还不舍得给能干的员工很高的工资，于是他规定：最高工资5000，所有人不可超过该工资。\r\n\r\n​	输入NANA和伙伴们的信息，请设计一个结构体储存这些信息，并设计一个函数模拟工作和涨工资过程，其参数是这样的结构体类型，返回同样的结构体类型，并输出员工信息。\r\n\r\n\r\n\r\n```plaintext\r\n输入格式\r\n第一行输入一个正整数 n，表示员工个数。\r\n第二行开始往下 n 行。每行首先是一个字符串表示员工姓名，再是一个整数表示员工年龄，再是一个整数为去年工资。\r\n\r\n输出格式\r\n输出 n 行，每行首先输出一个字符串表示姓名，再往后两个整数，表示经过一年的打工后大家的年龄和他们涨薪后的工资。以空格隔开。\r\n\r\n输入输出样例\r\n输入 \r\n3\r\nNANA 20 3600\r\nChenYu 23 3100\r\nVenti 2600 4200\r\n\r\n输出\r\nNANA 21 4500\r\nChenYu 24 3875\r\nVenti 2601 5000\r\n```\r\n代码如下：\r\n\r\n```c\r\n#include<bits/stdc++.h>\r\nusing namespace std;\r\nstruct Salary {\r\n	string name;\r\n	int age, salary;\r\n} a[100000];\r\n\r\nint main () {\r\n	int n;\r\n	cin >> n;\r\n	for (int i = 1; i <= n; i++) {\r\n		cin >> a[i].name >> a[i].age >> a[i].salary;\r\n		a[i].age++, a[i].salary = a[i].salary * 1.25;\r\n		if (a[i].salary > 5000) a[i].salary = 5000;\r\n		cout << a[i].name << \" \" << a[i].age << \" \" << a[i].salary << endl;\r\n	}\r\n	return 0;\r\n}\r\n```',1,1,349,'2024-08-22 01:31:24','2025-02-04 13:27:28',0,'lb6.jpg'),(77,'测试索引跳转','# 12\r\n4567541562146\r\n## 34\r\n4864896\r\n### 56\r\n4974896\r\n## 78\r\n278\r\n',12,1,22,'2025-02-04 11:07:10','2025-02-04 13:37:41',0,'31.jpg');
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
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'Coding-Study','代码学习','Venti-2.jpg',NULL),(2,'Project-Practice','项目实战','Venti-6.jpg',NULL),(3,'Plan','敬请期待','Venti-3.jpg',NULL),(4,'N-A-C-G','NACG','Venti-4.jpg',NULL),(5,'Travel','旅行日记','温迪-漫画封面.jpg',NULL),(12,'测试','AAA','Venti-5.jpg',NULL);
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
  PRIMARY KEY (`id`),
  KEY `user_comment_fk` (`user_id`),
  KEY `a_c_fk` (`article_id`),
  CONSTRAINT `c_a_fk` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_comment_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (20,4,1,'你还送','2024-05-23 05:39:36','2024-05-31 13:31:45'),(21,4,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 10:23:00','2024-05-31 13:32:01'),(22,29,1,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 10:54:13','2024-05-31 13:32:08'),(25,4,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:24:28','2024-05-31 13:32:15'),(27,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:52:58','2024-05-31 13:32:31'),(28,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:54:01','2024-05-31 13:32:34'),(29,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 12:20:22','2024-05-31 13:32:37'),(30,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 12:46:27','2024-05-31 13:32:41'),(31,29,2,'YRN','2024-05-24 13:35:47','2024-05-24 13:35:47'),(34,29,1,'123','2024-05-26 07:42:06','2024-05-26 07:42:06'),(35,29,2,'123','2024-05-26 07:48:34','2024-05-26 07:48:34'),(36,4,5,'123','2024-05-26 09:44:34','2024-05-26 09:44:34'),(37,29,2,'123','2024-05-26 10:41:01','2024-05-26 10:41:01'),(38,29,2,'456','2024-05-26 10:43:22','2024-05-26 10:43:22'),(39,29,2,'789','2024-05-26 10:43:31','2024-05-26 10:43:31'),(42,4,2,'喵喵喵','2024-05-27 15:36:11','2024-05-31 13:33:34'),(44,4,1,'喵喵喵','2024-05-29 06:58:35','2024-05-31 13:33:45'),(48,4,1,'喵喵喵','2024-05-29 10:43:18','2024-05-31 13:33:48'),(49,4,1,'喵喵喵','2024-05-29 10:43:19','2024-05-31 13:33:51'),(53,4,1,'耶耶耶耶耶','2024-05-29 10:53:49','2024-05-31 13:34:06'),(54,4,1,'c盘战士','2024-05-29 10:54:27','2024-05-31 13:34:24'),(56,4,2,'123','2024-05-29 11:22:51','2024-05-29 11:22:51'),(60,4,1,'完结撒花','2024-05-30 06:06:11','2024-05-31 13:34:40'),(69,4,1,'虚构史学家','2024-05-31 09:52:52','2024-05-31 13:35:44'),(72,13,1,'喵喵喵','2024-05-31 13:42:46','2024-05-31 13:42:46'),(73,16,1,'+10086','2024-05-31 13:47:50','2024-05-31 13:47:50'),(75,16,13,'阿巴阿巴','2024-06-01 04:31:52','2024-06-01 04:31:52'),(76,16,20,'小樱到此一游 请求卫宫士郎','2024-07-24 08:11:31','2024-07-24 08:11:31');
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES ('原神','https://ys.mihoyo.com/','当然是原神了!!','Venti-4.jpg',1),('崩坏-星穹铁道','https://sr.mihoyo.com/','我们在大量梗里发现了少量的游戏...','liuyingbangbu.jpg',2),('BiliBili','https://www.bilibili.com/','<小破站>','bilibili.jpg',3),('绝区零','https://zzz.mihoyo.com/main','豪玩！振刀爽！','ailian.jpg',6),('明日方舟','https://ak.hypergryph.com/','我看你才像植物大战僵尸','mingrifangzhou.png',13),('qqqqqq','qqqqqqqqqq','qqq','35.jpg',14),('qqqqqq','qqqqqqqqqq','qqq','35.jpg',15),('1234','123445','1431241542','屏幕截图 2024-06-01 224328.png',16);
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
) ENGINE=InnoDB AUTO_INCREMENT=67 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,1,'完结撒花！','2024-05-30 06:50:29'),(2,2,'完结撒花！','2024-05-30 06:50:29'),(4,5,'完结撒花！','2024-05-30 12:02:15'),(5,11,'完结撒花！','2024-05-30 12:02:15'),(6,6,'完结撒花！','2024-05-30 13:18:07'),(7,9,'完结撒花！','2024-05-30 13:18:07'),(10,9,'完结撒花！','2024-05-30 13:53:17'),(14,1,'完结撒花！','2024-05-31 10:04:02'),(15,11,'完结撒花！','2024-05-31 10:04:47'),(16,11,'完结撒花！','2024-05-31 10:08:47'),(18,1,'大搬家成功！','2024-05-31 13:28:36'),(20,13,'让我看看谁是第一个观众呢','2024-06-01 04:30:56'),(22,17,'123456','2024-06-10 03:34:36'),(25,1,'完结撒花！14','2024-05-30 06:50:29'),(26,1,'完结撒花！15','2024-05-30 06:50:29'),(27,1,'完结撒花！16','2024-05-30 06:50:29'),(28,1,'完结撒花！17','2024-05-30 06:50:29'),(29,1,'完结撒花！18','2024-05-30 06:50:29'),(36,19,'我是桀桀桀','2024-08-11 12:13:39'),(55,19,'啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 12:46:26'),(56,19,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 13:07:47'),(57,19,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 13:08:39'),(58,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:28'),(59,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:43'),(60,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:43'),(62,1,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-21 13:49:34'),(63,1,'test_0921','2024-09-21 02:13:42'),(64,1,'test_0921','2024-09-21 02:31:18'),(65,1,'test_0921','2024-09-21 02:32:24'),(66,1,'1111111','2025-02-03 10:56:49');
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'CLN','mH7bNTkcrZuPAA==','2024-05-18 12:49:38','A','35.jpg',0,'changbingmushao@qq.com'),(2,'YRN','821212yrn','2024-05-23 08:17:02','A','yrn.jpg',0,NULL),(5,'III','456789iii','2024-05-23 09:16:50','B','31.jpg',0,NULL),(6,'FFF','111111fff','2024-05-23 12:21:25','B','34.jpg',0,NULL),(9,'zzx','123456dzs','2024-05-27 08:29:59','B','36.jpg',0,NULL),(11,'倒置头','222222eee','2024-05-30 04:05:34','B','daozhi.jpg',0,NULL),(13,'gyh','20051012gyh','2024-06-01 04:24:10','A','mrfz.jpg',0,NULL),(15,'1','1','2024-06-02 09:08:39','B','boli.jpg',0,NULL),(16,'111','123456aaa','2024-06-02 09:09:38','B','onlyBlack.jpg',0,NULL),(17,'CLD','124','2024-06-10 03:33:06','B','renma.jpg',0,NULL),(18,'库洛米','250420cln','2024-07-24 03:28:52','A','kuluomi.jpg',0,NULL),(19,'桀桀桀','123456jjj','2024-07-24 04:41:35','B','mojinghaha.jpg',0,'hahaha@qq.com'),(20,'樱','123456yyy','2024-07-24 08:09:01','B','moren.jpg',0,NULL),(34,'小黑','srHkWMcoGOPwtA==','2024-10-18 10:11:21','B','屏幕截图 2024-10-16 101352.png',0,'juewangwang@qq.com'),(35,'JWW','srHkWMcoGOPwtA==','2024-10-18 10:15:44','B','屏幕截图 2024-10-16 221438.png',0,'48588@qq.com'),(36,'AAA','iRKexQt7QFeoLw==','2025-01-23 08:21:54','B','屏幕截图 2024-05-28 201303.png',0,''),(38,'FuHua','SXeR07EdIl1PWQ==','2025-01-23 10:10:11','B','37.jpg',0,'');
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

-- Dump completed on 2025-02-04 21:39:40
