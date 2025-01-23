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
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (4,'今天玩什么？','# 当然是致命公司了!\n#### 但是泰坦图第1e9+7次满货死亡.这种游戏不连麦玩，自己挂在远处队友都不知道，如果侥幸活下来，队友直接开飞船跑了，这种情况 我真的是很绝望啊',4,1,72,'2024-05-18 12:51:04','2024-09-20 13:19:27',0,'zmgs.png'),(13,'cors使用','#### corsConfig := cors.DefaultConfig()\n#### corsConfig.AllowOrigins = []string{\"*\"}\n#### corsConfig.AllowMethods = []string{\"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"OPTIONS\"}\n#### r.Use(cors.New(corsConfig))\n',1,1,14,'2024-05-19 11:03:28','2025-01-23 09:21:40',0,'cors.png'),(16,'评论水区','直接发吧\n',12,1,25,'2024-05-19 12:20:38','2024-08-12 07:14:32',0,'Venti-5.jpg'),(29,'js gin链接','fetch 正常写       实时反馈控制台   就这样',1,1,6,'2024-05-21 12:55:09','2024-10-18 10:27:45',0,'34.jpg'),(73,'小游戏','# 一些小游戏\r\n### https://github.com/shaozicln/Game<div style=\"color: rgba(0,0,0,0.6); font-size: 30px; text-align: center;\"><br><br><br>----- 有关Project - Practice的一切都在 -----<br>https://github.com/shaozicln</div>',2,1,1,'2025-01-23 09:52:49','2025-01-23 09:52:53',0,'键盘.jpeg');
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
INSERT INTO `categories` VALUES (1,'Coding-Study','代码学习','Venti-2.jpg',NULL),(2,'Project-Practice','项目实战','Venti-6.jpg',NULL),(3,'Life','生活趣事','Venti-3.jpg',NULL),(4,'N-A-C-G','NACG','Venti-4.jpg',NULL),(12,'测试','AAA','Venti-5.jpg',NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,1,'完结撒花！','2024-05-30 06:50:29'),(2,2,'完结撒花！','2024-05-30 06:50:29'),(4,5,'完结撒花！','2024-05-30 12:02:15'),(5,11,'完结撒花！','2024-05-30 12:02:15'),(6,6,'完结撒花！','2024-05-30 13:18:07'),(7,9,'完结撒花！','2024-05-30 13:18:07'),(10,9,'完结撒花！','2024-05-30 13:53:17'),(14,1,'完结撒花！','2024-05-31 10:04:02'),(15,11,'完结撒花！','2024-05-31 10:04:47'),(16,11,'完结撒花！','2024-05-31 10:08:47'),(18,1,'大搬家成功！','2024-05-31 13:28:36'),(20,13,'让我看看谁是第一个观众呢','2024-06-01 04:30:56'),(22,17,'123456','2024-06-10 03:34:36'),(25,1,'完结撒花！14','2024-05-30 06:50:29'),(26,1,'完结撒花！15','2024-05-30 06:50:29'),(27,1,'完结撒花！16','2024-05-30 06:50:29'),(28,1,'完结撒花！17','2024-05-30 06:50:29'),(29,1,'完结撒花！18','2024-05-30 06:50:29'),(36,19,'我是桀桀桀','2024-08-11 12:13:39'),(55,19,'啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 12:46:26'),(56,19,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 13:07:47'),(57,19,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-11 13:08:39'),(58,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:28'),(59,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:43'),(60,19,'爱啊爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱爱','2024-08-11 13:09:43'),(62,1,'啊啊啊啊啊啊啊啊啊啊啊啊啊啊','2024-08-21 13:49:34'),(63,1,'test_0921','2024-09-21 02:13:42'),(64,1,'test_0921','2024-09-21 02:31:18'),(65,1,'test_0921','2024-09-21 02:32:24');
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
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'CLN','mH7bNTkcrZuPAA==','2024-05-18 12:49:38','A','35.jpg',0,'changbingmushao@qq.com'),(2,'YRN','821212yrn','2024-05-23 08:17:02','A','yrn.jpg',0,NULL),(5,'III','456789iii','2024-05-23 09:16:50','B','31.jpg',0,NULL),(6,'FFF','111111fff','2024-05-23 12:21:25','B','34.jpg',0,NULL),(9,'zzx','123456dzs','2024-05-27 08:29:59','B','36.jpg',0,NULL),(11,'倒置头','222222eee','2024-05-30 04:05:34','B','daozhi.jpg',0,NULL),(13,'gyh','20051012gyh','2024-06-01 04:24:10','A','mrfz.jpg',0,NULL),(15,'1','1','2024-06-02 09:08:39','B','boli.jpg',0,NULL),(16,'111','123456aaa','2024-06-02 09:09:38','B','onlyBlack.jpg',0,NULL),(17,'CLD','124','2024-06-10 03:33:06','B','renma.jpg',0,NULL),(18,'库洛米','250420cln','2024-07-24 03:28:52','A','kuluomi.jpg',0,NULL),(19,'桀桀桀','123456jjj','2024-07-24 04:41:35','B','mojinghaha.jpg',0,'hahaha@qq.com'),(20,'樱','123456yyy','2024-07-24 08:09:01','B','moren.jpg',0,NULL),(34,'小黑','srHkWMcoGOPwtA==','2024-10-18 10:11:21','B','屏幕截图 2024-10-16 101352.png',0,'juewangwang@qq.com'),(35,'JWW','srHkWMcoGOPwtA==','2024-10-18 10:15:44','B','屏幕截图 2024-10-16 221438.png',0,'48588@qq.com'),(36,'AAA','iRKexQt7QFeoLw==','2025-01-23 08:21:54','B','屏幕截图 2024-05-28 201303.png',0,'');
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

-- Dump completed on 2025-01-23 17:56:08
