-- MySQL dump 10.16  Distrib 10.3.1-MariaDB, for Win64 (AMD64)
--
-- Host: 192.168.10.190    Database: mariadb
-- ------------------------------------------------------
-- Server version	10.0.29-MariaDB-0ubuntu0.16.10.1

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

CREATE USER 'mediadb' IDENTIFIED BY 'mediaDB';
CREATE USER 'mediadb'@'localhost' IDENTIFIED BY 'mediaDB';


DROP DATABASE IF EXISTS `mediadb`;

CREATE DATABASE `mediadb` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

use mediadb;

GRANT ALL PRIVILEGES ON mediadb.* TO 'mediadb'@'%';
GRANT ALL PRIVILEGES ON mediadb.* TO 'mediadb'@'localhost';
FLUSH PRIVILEGES;


SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `filelist` CASCADE;

CREATE TABLE `filelist` (
  `fileid` int(11) NOT NULL AUTO_INCREMENT,
  `filename` varchar(254) NOT NULL,
  `filesuffix` varchar(5) DEFAULT NULL,
  `filelocation` varchar(254) DEFAULT NULL,
  `filesize` float DEFAULT NULL,
  `filehash` varchar(254) DEFAULT NULL,
  `filedate` datetime NOT NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `dupecount` int(11) DEFAULT NULL,
  PRIMARY KEY (`fileid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='list of files on the drive';

DROP TABLE IF EXISTS `targetinfo` CASCADE;

CREATE TABLE `targetinfo` (
  `targetid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11) NOT NULL,
  `filename` varchar(254) NOT NULL,
  `filesuffix` varchar(5) DEFAULT NULL,
  `filelocation` varchar(254) DEFAULT NULL,
  `filehash` varchar(254) DEFAULT NULL,
  `filesize` float DEFAULT NULL,
  `filedate` datetime NOT NULL,
  `fileaction` char(2) NOT NULL,
  `validated` tinyint(1) NOT NULL DEFAULT '0',
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`targetid`),
  KEY `idx_targetinfo` (`fileid`),
  CONSTRAINT `fk_targetinfo_filelist` FOREIGN KEY (`fileid`) REFERENCES `filelist` (`fileid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Target information for the files to be moved or copied';

DROP TABLE IF EXISTS `dupes` CASCADE;

CREATE TABLE `dupes` (
  `dupeid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11) NOT NULL,
  `dupefileid` int(11) NOT NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`dupeid`),
  KEY `idx_dupes` (`fileid`),
  CONSTRAINT `fk_dupes_filelist` FOREIGN KEY (`fileid`) REFERENCES `filelist` (`fileid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=' flag any duplicate files';

DROP TABLE IF EXISTS `exifinfo` CASCADE;

CREATE TABLE `exifinfo` (
  `exifid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11) NOT NULL,
  `accessdate` datetime NOT NULL,
  `comments` varchar(254) DEFAULT NULL,
  `createdate` datetime NOT NULL,
  `facecoords` varchar(254) DEFAULT NULL,
  `faces` varchar(254) DEFAULT NULL,
  `gpsaltitude` varchar(254) DEFAULT NULL,
  `gpsaltituderef` varchar(254) DEFAULT NULL,
  `gpsdatetime` datetime,
  `gpslatitude` varchar(254) DEFAULT NULL,
  `gpslatituderef` varchar(254) DEFAULT NULL,
  `gpslongitude` varchar(254) DEFAULT NULL,
  `gpslongituderef` varchar(254) DEFAULT NULL,
  `gpsmapdatum` varchar(254) DEFAULT NULL,
  `gpsprocessingmethod` varchar(254) DEFAULT NULL,
  `gpsversionid` varchar(254) DEFAULT NULL,
  `imageDescription` varchar(254) DEFAULT NULL,
  `make` varchar(254) DEFAULT NULL,
  `model` varchar(254) DEFAULT NULL,
  `modifydate` datetime NOT NULL,
  `orientation` varchar(254) DEFAULT NULL,
  `tags` varchar(254) DEFAULT NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`exifid`),
  KEY `idx_exifinfo` (`fileid`),
  CONSTRAINT `fk_exifinfo_filelist` FOREIGN KEY (`fileid`) REFERENCES `filelist` (`fileid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='EXIF Info of the file';

DROP TABLE IF EXISTS `thumbnails` CASCADE;

CREATE TABLE `thumbnails` (
  `thid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11) NOT NULL,
  `targetid` int(11) NOT NULL,
  `filename` varchar(254) NOT NULL,
  `filesuffix` varchar(254) DEFAULT NULL,
  `filesize` varchar(254) DEFAULT NULL,
  `thsize` varchar(254) NOT NULL,
  `fileloc` varchar(254) DEFAULT NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`thid`),
  KEY `idx_thumbnails` (`targetid`),
  KEY `idx1_thumbnails` (`fileid`),
  CONSTRAINT `fk_thumbnails_filelist` FOREIGN KEY (`fileid`) REFERENCES `filelist` (`fileid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_thumbnails_targetinfo` FOREIGN KEY (`targetid`) REFERENCES `targetinfo` (`targetid`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `regionInfo` CASCADE;

CREATE TABLE `regionInfo` (
  `regionid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11),
  `filehash` varchar(254) DEFAULT NULL,
  `name` varchar(254) NOT NULL,
  `type` varchar(5) DEFAULT NULL,
  `areaH` float DEFAULT NULL,
  `areaW` float DEFAULT NULL,
  `areaX` float DEFAULT NULL,
  `areaY` float DEFAULT NULL,  
  `areaUnit` varchar(254) DEFAULT NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`regionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='list of regions(face) on an image/file';

DROP TABLE IF EXISTS `imageInfo` CASCADE;

CREATE TABLE `imageInfo` (
  `imgid` int(11) NOT NULL AUTO_INCREMENT,
  `fileid` int(11),
  `filehash` varchar(254) DEFAULT NULL,
  `dimH` int DEFAULT NULL,
  `dimW` int DEFAULT NULL,
  `dimUnit` varchar(254) DEFAULT NULL,
  `rotation` tinyint unsigned default NULL,
  `latitude` double(8, 6) default NULL,
  `longitude` double(9, 6) default NULL,
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`imgid`),
  KEY `images_i2` (`fileid`),
  KEY `images_i3` (`filehash`),
  KEY `images_i4` (`dimH`),
  KEY `images_i5` (`dimW`),
  KEY `images_i1` (`rowaction`),
  KEY `images_i6` (`latitude`),
  KEY `images_i7` (`longitude`),
  KEY `lastmodified` (`rowactiondatetime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='informaton about the image on the file, dimensions etc. The location of faces is relative to these dimensions';

DROP TABLE IF EXISTS `sessionInfo` CASCADE;

CREATE TABLE `sessionInfo` (
  `sessionid` int(11) NOT NULL AUTO_INCREMENT,
  `starttime` datetime,
  `stoptime`  datetime,
  `totalFiles` int,
  `totalBytes`  int,
  `examinedFiles` int,
  `examinedBytes` int,
  `uniqueFiles`   int,
  `uniqueBytes`   int,
  `runTime`   float,
  `logFile`   varchar(254),
  `rowaction` char(2) NOT NULL,
  `rowactiondatetime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`sessionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='informaton about the image on the file, dimensions etc. The location of faces is relative to these dimensions';

SET FOREIGN_KEY_CHECKS = 1;