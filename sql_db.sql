-- MySQL dump 10.14  Distrib 5.5.41-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: cdr
-- ------------------------------------------------------
-- Server version	5.5.41-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `agents`
--

DROP TABLE IF EXISTS `recordingdata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `recordingdata` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `calldate` varchar(255) NOT NULL,
  `clid` varchar(255) NOT NULL UNIQUE,
  `src` varchar(255) NOT NULL,
  `dst` varchar(255) NOT NULL,
  `duration` varchar(255) NOT NULL,
  `billsec` varchar(255) NOT NULL,
  `disposition` varchar(255) NOT NULL,
  `accountcode` varchar(255) NOT NULL,
  `uniqueid` varchar(255) NOT NULL,
  `did` varchar(255) NOT NULL,
  `recordingfile` varchar(255) NOT NULL,
  `diskfilepath` varchar(255) NOT NULL,
  `s3fileurl` varchar(255) NOT NULL,
  `office` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

/*DROP TABLE IF EXISTS `credentials`;
40101 SET @saved_cs_client     = @@character_set_client ;
40101 SET character_set_client = utf8 ;

CREATE TABLE `credentials` (
	  `id` int(11) NOT NULL AUTO_INCREMENT,
	  `title` varchar(255) DEFAULT NULL,
	  `username` varchar(255) DEFAULT NULL,
	  `description` varchar(255) DEFAULT NULL,
	  `tag` varchar(255) DEFAULT NULL,
	  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=230 DEFAULT CHARSET=latin1;
*/

/*DROP TABLE IF EXISTS `users`;*/
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
/*
CREATE TABLE `users` (
	  `id` int(11) NOT NULL AUTO_INCREMENT,
	  `username` varchar(255) DEFAULT NULL,
	  `email` varchar(255) DEFAULT NULL,
	  `group_tag` varchar(255) DEFAULT NULL,
	  `password` varchar(255) DEFAULT NULL,
	  `active` TINYINT(1) NOT NULL,
	  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=230 DEFAULT CHARSET=latin1;
*/

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cdrs`
--
/*
DROP TABLE IF EXISTS `cdrs`;


CREATE TABLE `cdrs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `accountcode` varchar(255) DEFAULT NULL,
  `source` varchar(255) DEFAULT NULL,
  `dest` varchar(255) DEFAULT NULL,
  `dcontext` varchar(255) DEFAULT NULL,
  `callerid` varchar(255) DEFAULT NULL,
  `channel` varchar(255) DEFAULT NULL,
  `dstchannel` varchar(255) DEFAULT NULL,
  `lastapp` varchar(255) DEFAULT NULL,
  `lastdata` varchar(255) DEFAULT NULL,
  `start` varchar(255) DEFAULT NULL,
  `answer` varchar(255) DEFAULT NULL,
  `end` varchar(255) DEFAULT NULL,
  `duration` varchar(255) DEFAULT NULL,
  `billsec` varchar(255) DEFAULT NULL,
  `disposition` varchar(255) DEFAULT NULL,
  `amaflags` varchar(255) DEFAULT NULL,
  `uniqueid` varchar(255) NOT NULL,
  `userfield` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqueid` (`uniqueid`)
) ENGINE=InnoDB AUTO_INCREMENT=139480 DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `desks`;


CREATE TABLE `desks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `accountcode` varchar(255) DEFAULT NULL,
  `source` varchar(255) DEFAULT NULL,
  `dest` varchar(255) DEFAULT NULL,
  `dcontext` varchar(255) DEFAULT NULL,
  `callerid` varchar(255) DEFAULT NULL,
  `channel` varchar(255) DEFAULT NULL,
  `dstchannel` varchar(255) DEFAULT NULL,
  `lastapp` varchar(255) DEFAULT NULL,
  `lastdata` varchar(255) DEFAULT NULL,
  `start` varchar(255) DEFAULT NULL,
  `answer` varchar(255) DEFAULT NULL,
  `end` varchar(255) DEFAULT NULL,
  `duration` varchar(255) DEFAULT NULL,
  `billsec` varchar(255) DEFAULT NULL,
  `disposition` varchar(255) DEFAULT NULL,
  `amaflags` varchar(255) DEFAULT NULL,
  `uniqueid` varchar(255) DEFAULT NULL,
  `userfield` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/

