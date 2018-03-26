-- MySQL dump 10.16  Distrib 10.1.12-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: catkeeper
-- ------------------------------------------------------
-- Server version	10.1.12-MariaDB

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
-- Table structure for table `physicalmachine`
--

DROP TABLE IF EXISTS `physicalmachine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `physicalmachine` (
  `ipaddress` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `cputype` varchar(60) DEFAULT NULL,
  `cpu` int(10) DEFAULT NULL,
  `cpufreq` int(10) DEFAULT NULL,
  `cpusocket` int(10) DEFAULT NULL,
  `cpukernel` int(10) DEFAULT NULL,
  `cputhread` int(10) DEFAULT NULL,
  `numa` int(10) DEFAULT NULL,
  `memory` int(12) DEFAULT NULL,
  PRIMARY KEY (`ipaddress`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `virtualmachine`
--

DROP TABLE IF EXISTS `virtualmachine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `virtualmachine` (
  `uuid` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `hostipaddress` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `attachments` varchar(255) DEFAULT NULL,
  `cpu` int(11) DEFAULT NULL,
  `mem` int(11) DEFAULT NULL,
  `disk` int(11) DEFAULT NULL,
  `imagename` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `sysdisk` varchar(255) DEFAULT NULL,
  `user` varchar(60) DEFAULT NULL,
  `createtime` varchar(60) DEFAULT NULL,
  `updatetime` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `volume`
--

DROP TABLE IF EXISTS `volume`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `volume` (
  `uuid` varchar(255) NOT NULL,
  `size` int(11) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `volumetype` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `createat` varchar(255) DEFAULT NULL,
  `updateat` varchar(255) DEFAULT NULL,
  `attachments` varchar(255) DEFAULT NULL,
  `hostip` varchar(255) DEFAULT NULL,
  `datatype` varchar(255) DEFAULT NULL,
  `user` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

CREATE TABLE `physicalmachinediskinfo` (
  `uuid` varchar(255) DEFAULT NULL,
  `Ip` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `capility` int(12) DEFAULT NULL,
  `status` varchar(60) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `virtualmachinelog` (
  `uuid` varchar(255) CHARACTER SET latin1 NOT NULL,
  `vmuuid` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `action` varchar(50) CHARACTER SET latin1 DEFAULT NULL,
  `status` varchar(50) CHARACTER SET latin1 DEFAULT NULL,
  `operatetime` varchar(255) CHARACTER SET latin1 DEFAULT NULL,
  `describ` varchar(255) DEFAULT NULL,
  `user` varchar(60) DEFAULT NULL,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Dump completed on 2018-03-07 20:34:41
