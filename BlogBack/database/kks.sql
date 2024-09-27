-- MySQL dump 10.13  Distrib 8.0.37, for Win64 (x86_64)
--
-- Host: localhost    Database: kh
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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applications`
--

LOCK TABLES `applications` WRITE;
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` VALUES (1,'桀桀桀','hahaha@qq.com','桀桀桀の网站','https://www.aaa.com/','来自桀桀桀的aaa网站','Venti-6.jpg'),(3,'桀桀桀','hahaha@qq.com','111111111','111111111111','11111111111111111','35.jpg'),(5,'CLN','changbingmushao@qq.com','qqqqqq','qqqqqqqqqq','qqq','35.jpg');
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
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (4,'今天玩什么？','# 当然是致命公司了!\n#### 但是泰坦图第1e9+7次满货死亡.这种游戏不连麦玩，自己挂在远处队友都不知道，如果侥幸活下来，队友直接开飞船跑了，这种情况 我真的是很绝望啊',4,1,72,'2024-05-18 12:51:04','2024-09-20 13:19:27',0,'zmgs.png'),(13,'cors使用','#### corsConfig := cors.DefaultConfig()\n#### corsConfig.AllowOrigins = []string{\"*\"}\n#### corsConfig.AllowMethods = []string{\"GET\", \"POST\", \"PUT\", \"PATCH\", \"DELETE\", \"OPTIONS\"}\n#### r.Use(cors.New(corsConfig))\n',2,1,14,'2024-05-19 11:03:28','2024-08-12 07:37:50',0,'cors.png'),(16,'评论水区','直接发吧\n',12,1,25,'2024-05-19 12:20:38','2024-08-12 07:14:32',0,'Venti-5.jpg'),(29,'js gin链接','fetch 正常写       实时反馈控制台   就这样',1,1,5,'2024-05-21 12:55:09','2024-08-12 05:49:27',0,'34.jpg'),(32,'《荷塘月色》——朱自清','这几天心里颇不宁静。今晚在院子里坐着乘凉，忽然想起日日走过的荷塘，在这满月的光里，总该另有一番样子吧。月亮渐渐地升高了，墙外马路上孩子们的欢笑，已经听不见了；妻在屋里拍着闰儿，迷迷糊糊地哼着眠歌。我悄悄地披了大衫，带上门出去。\n　　沿着荷塘，是一条曲折的小煤屑路。这是一条幽僻的路；白天也少人走，夜晚更加寂寞。荷塘四面，长着许多树，蓊蓊郁郁的。路的一旁，是些杨柳，和一些不知道名字的树。没有月光的晚上，这路上阴森森的，有些怕人。今晚却很好，虽然月光也还是淡淡的。\n　　路上只我一个人，背着手踱着。这一片天地好像是我的；我也像超出了平常的自己，到了另一世界里。我爱热闹，也爱冷静；爱群居，也爱独处。像今晚上，一个人在这苍茫的月下，什么都可以想，什么都可以不想，便觉是个自由的人。白天里一定要做的事，一定要说的话，现在都可不理。这是独处的妙处，我且受用这无边的荷香月色好了。\n　　曲曲折折的荷塘上面，弥望的是田田的叶子。叶子出水很高，像亭亭的舞女的裙。层层的叶子中间，零星地点缀着些白花，有袅娜地开着的，有羞涩地打着朵儿的；正如一粒粒的明珠，又如碧天里的星星，又如刚出浴的美人。微风过处，送来缕缕清香，仿佛远处高楼上渺茫的歌声似的。这时候叶子与花也有一丝的颤动，像闪电般，霎时传过荷塘的那边去了。叶子本是肩并肩密密地挨着，这便宛然有了一道凝碧的波痕。叶子底下是脉脉的流水，遮住了，不能见一些颜色；而叶子却更见风致了。\n　　月光如流水一般，静静地泻在这一片叶子和花上。薄薄的青雾浮起在荷塘里。叶子和花仿佛在牛乳中洗过一样；又像笼着轻纱的梦。虽然是满月，天上却有一层淡淡的云，所以不能朗照；但我以为这恰是到了好处——酣眠固不可少，小睡也别有风味的。月光是隔了树照过来的，高处丛生的灌木，落下参差的斑驳的黑影，峭楞楞如鬼一般；弯弯的杨柳的稀疏的倩影，却又像是画在荷叶上。塘中的月色并不均匀；但光与影有着和谐的旋律，如梵婀玲上奏着的名曲。\n　　荷塘的四面，远远近近，高高低低都是树，而杨柳最多。这些树将一片荷塘重重围住；只在小路一旁，漏着几段空隙，像是特为月光留下的。树色一例是阴阴的，乍看像一团烟雾；但杨柳的丰姿，便在烟雾里也辨得出。树梢上隐隐约约的是一带远山，只有些大意罢了。树缝里也漏着一两点路灯光，没精打采的，是渴睡人的眼。这时候最热闹的，要数树上的蝉声与水里的蛙声；但热闹是它们的，我什么也没有。\n　　忽然想起采莲的事情来了。采莲是江南的旧俗，似乎很早就有，而六朝时为盛；从诗歌里可以约略知道。采莲的是少年的女子，她们是荡着小船，唱着艳歌去的。采莲人不用说很多，还有看采莲的人。那是一个热闹的季节，也是一个风流的季节。梁元帝《采莲赋》里说得好：\n　　于是妖童媛女，荡舟心许；鷁首徐回，兼传羽杯；欋将移而藻挂，船欲动而萍开。尔其纤腰束素，迁延顾步；夏始春余，叶嫩花初，恐沾裳而浅笑，畏倾船而敛裾。\n　　可见当时嬉游的光景了。这真是有趣的事，可惜我们现在早已无福消受了。\n　　于是又记起《西洲曲》里的句子：\n　　采莲南塘秋，莲花过人头；低头弄莲子，莲子清如水。今晚若有采莲人，这儿的莲花也算得“过人头”了；只不见一些流水的影子，是不行的。这令我到底惦着江南了。——这样想着，猛一抬头，不觉已是自己的门前；轻轻地推门进去，什么声息也没有，妻已睡熟好久了。\n　\n　　1927年7月，北京清华园。',4,1,13,'2024-05-28 13:18:55','2024-08-12 05:57:56',0,'htys.jpg'),(38,'重合的开局 不同的震撼','《终结的炽天使》和《反叛的鲁路修》',4,1,12,'2024-06-01 10:45:00','2024-08-12 07:20:55',0,'luluxiu.jpg'),(40,'生活记录','aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa',3,1,2,'2024-08-09 07:23:51','2024-08-12 05:35:00',0,'Venti-6.jpg'),(49,'测试文章0812','# 一级标题测试',12,1,29,'2024-08-12 01:11:20','2024-08-12 07:20:03',0,'36.jpg'),(50,'Web测试文章0812','# 二级标题\n## 三级标题\n![图片测试](../Public/Pictures/Venti-1.jpg)',12,1,21,'2024-08-12 01:17:35','2024-08-12 08:30:26',0,'37.jpg'),(51,'测试文章2222222','### 气温起伏企鹅',12,1,9,'2024-08-12 03:15:35','2024-08-16 02:30:41',0,'cln-left.jpg'),(53,'1243','chuoujjhdbbss fWA FW f\n',1,1,4,'2024-08-16 02:40:21','2024-09-20 13:19:20',0,'Venti-6');
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
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'Front-end development','前端内容','Venti-2.jpg',NULL),(2,'Back-end development','后端内容','Venti-6.jpg',NULL),(3,'Career planning','职业想法','Venti-3.jpg',NULL),(4,'N-A-C-G','NACG','Venti-4.jpg',NULL),(12,'测试','AAA','Venti-5.jpg',NULL);
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
INSERT INTO `comments` VALUES (20,4,1,'你还送','2024-05-23 05:39:36','2024-05-31 13:31:45'),(21,4,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 10:23:00','2024-05-31 13:32:01'),(22,29,1,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 10:54:13','2024-05-31 13:32:08'),(25,4,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:24:28','2024-05-31 13:32:15'),(27,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:52:58','2024-05-31 13:32:31'),(28,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 11:54:01','2024-05-31 13:32:34'),(29,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 12:20:22','2024-05-31 13:32:37'),(30,29,2,'哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈','2024-05-24 12:46:27','2024-05-31 13:32:41'),(31,29,2,'YRN','2024-05-24 13:35:47','2024-05-24 13:35:47'),(34,29,1,'123','2024-05-26 07:42:06','2024-05-26 07:42:06'),(35,29,2,'123','2024-05-26 07:48:34','2024-05-26 07:48:34'),(36,4,5,'123','2024-05-26 09:44:34','2024-05-26 09:44:34'),(37,29,2,'123','2024-05-26 10:41:01','2024-05-26 10:41:01'),(38,29,2,'456','2024-05-26 10:43:22','2024-05-26 10:43:22'),(39,29,2,'789','2024-05-26 10:43:31','2024-05-26 10:43:31'),(42,4,2,'喵喵喵','2024-05-27 15:36:11','2024-05-31 13:33:34'),(44,4,1,'喵喵喵','2024-05-29 06:58:35','2024-05-31 13:33:45'),(48,4,1,'喵喵喵','2024-05-29 10:43:18','2024-05-31 13:33:48'),(49,4,1,'喵喵喵','2024-05-29 10:43:19','2024-05-31 13:33:51'),(53,4,1,'耶耶耶耶耶','2024-05-29 10:53:49','2024-05-31 13:34:06'),(54,4,1,'c盘战士','2024-05-29 10:54:27','2024-05-31 13:34:24'),(56,4,2,'123','2024-05-29 11:22:51','2024-05-29 11:22:51'),(60,4,1,'完结撒花','2024-05-30 06:06:11','2024-05-31 13:34:40'),(64,32,11,'朱自清先生写得真是太棒了','2024-05-30 07:28:35','2024-05-31 13:35:00'),(69,4,1,'虚构史学家','2024-05-31 09:52:52','2024-05-31 13:35:44'),(71,32,1,'朱自清先生写得真是太棒了','2024-05-31 10:01:19','2024-05-31 13:35:58'),(72,13,1,'喵喵喵','2024-05-31 13:42:46','2024-05-31 13:42:46'),(73,16,1,'+10086','2024-05-31 13:47:50','2024-05-31 13:47:50'),(75,16,13,'阿巴阿巴','2024-06-01 04:31:52','2024-06-01 04:31:52'),(76,16,20,'小樱到此一游 请求卫宫士郎','2024-07-24 08:11:31','2024-07-24 08:11:31');
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
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES ('原神','https://ys.mihoyo.com/','当然是原神了!!','Venti-4.jpg',1),('崩坏-星穹铁道','https://sr.mihoyo.com/','我们在大量梗里发现了少量的游戏...','liuyingbangbu.jpg',2),('BiliBili','https://www.bilibili.com/','<小破站>','bilibili.jpg',3),('绝区零','https://zzz.mihoyo.com/main','豪玩！振刀爽！','ailian.jpg',6),('明日方舟','https://ak.hypergryph.com/','我看你才像植物大战僵尸','mingrifangzhou.png',13),('qqqqqq','qqqqqqqqqq','qqq','35.jpg',14),('qqqqqq','qqqqqqqqqq','qqq','35.jpg',15);
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
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'CLN','250420cln','2024-05-18 12:49:38','A','35.jpg',0,'changbingmushao@qq.com'),(2,'YRN','821212yrn','2024-05-23 08:17:02','A','yrn.jpg',0,NULL),(5,'III','456789iii','2024-05-23 09:16:50','B','31.jpg',0,NULL),(6,'FFF','111111fff','2024-05-23 12:21:25','B','34.jpg',0,NULL),(9,'zzx','123456dzs','2024-05-27 08:29:59','B','36.jpg',0,NULL),(11,'倒置头','222222eee','2024-05-30 04:05:34','B','daozhi.jpg',0,NULL),(13,'gyh','20051012gyh','2024-06-01 04:24:10','A','mrfz.jpg',0,NULL),(15,'1','1','2024-06-02 09:08:39','B','boli.jpg',0,NULL),(16,'111','123456aaa','2024-06-02 09:09:38','B','onlyBlack.jpg',0,NULL),(17,'CLD','124','2024-06-10 03:33:06','B','renma.jpg',0,NULL),(18,'库洛米','250420cln','2024-07-24 03:28:52','A','kuluomi.jpg',0,NULL),(19,'桀桀桀','123456jjj','2024-07-24 04:41:35','B','mojinghaha.jpg',0,'hahaha@qq.com'),(20,'樱','123456yyy','2024-07-24 08:09:01','B','moren.jpg',0,NULL);
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

-- Dump completed on 2024-09-26 14:01:43
