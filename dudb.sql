# ************************************************************
# Sequel Pro SQL dump
# Version 4135
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.5.34)
# Database: dudb
# Generation Time: 2015-07-09 02:02:43 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table t_apps
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_apps`;

CREATE TABLE `t_apps` (
  `app_id` bigint(20) NOT NULL,
  `app_key` varchar(255) DEFAULT NULL,
  `dev_id` bigint(20) DEFAULT NULL,
  `cs_logourl` varchar(255) NOT NULL DEFAULT '',
  `cs_name` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `ind_applist_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_bind_phone
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_bind_phone`;

CREATE TABLE `t_bind_phone` (
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `phonenum` varchar(18) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '用户绑定的手机号',
  `bind_date` datetime NOT NULL COMMENT '手机号绑定的日期时间',
  PRIMARY KEY (`phonenum`),
  KEY `idx_uid` (`uid`),
  CONSTRAINT `uid_fk` FOREIGN KEY (`uid`) REFERENCES `t_user_info` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户绑定的手机号，一个用户ID可以绑定多个手机号';



# Dump of table t_blacklist
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_blacklist`;

CREATE TABLE `t_blacklist` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `fuid` bigint(20) NOT NULL DEFAULT '0',
  `itime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` int(1) NOT NULL DEFAULT '0' COMMENT '0 ï¼šæ­£å¸¸ï¼Œ1ï¼šå·²åˆ é™¤',
  `last_modify_date` int(14) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_fuid` (`uid`,`fuid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table t_csmsg_history
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_csmsg_history`;

CREATE TABLE `t_csmsg_history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `fromcid` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `tocid` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `msg` varchar(512) NOT NULL DEFAULT '',
  `itime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `msgid` bigint(20) NOT NULL DEFAULT '0',
  `appkey` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `index_appkey` (`appkey`),
  KEY `index_fromcid` (`fromcid`),
  KEY `index_tocid` (`tocid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_customservice
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_customservice`;

CREATE TABLE `t_customservice` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `account` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `password` varchar(33) NOT NULL DEFAULT '',
  `appkey` varchar(64) NOT NULL DEFAULT '',
  `nick_name` varchar(64) NOT NULL DEFAULT '',
  `image_id` varchar(64) NOT NULL DEFAULT '',
  `email` varchar(33) NOT NULL DEFAULT '',
  `tel` varchar(33) NOT NULL DEFAULT '',
  `enable` int(4) NOT NULL DEFAULT '1',
  `reg_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del` int(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_2` (`uid`),
  UNIQUE KEY `appkey_account` (`account`,`appkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_developer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_developer`;

CREATE TABLE `t_developer` (
  `did` int(14) unsigned NOT NULL AUTO_INCREMENT,
  `dname` varchar(64) NOT NULL DEFAULT '',
  `dpassword` varchar(64) NOT NULL DEFAULT '',
  `dprofile` varchar(512) NOT NULL DEFAULT '',
  PRIMARY KEY (`did`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_devicetoken
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_devicetoken`;

CREATE TABLE `t_devicetoken` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `token` varchar(128) DEFAULT NULL,
  `open_flag` int(4) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_index` (`uid`),
  KEY `Iindex_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_friendship
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_friendship`;

CREATE TABLE `t_friendship` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid1` bigint(11) unsigned DEFAULT NULL,
  `uid2` bigint(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table t_kkidpool
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_kkidpool`;

CREATE TABLE `t_kkidpool` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `kkid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'å¾…æ”¾kkå·',
  `out_flag` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0-æœªæ”¾å‡º 1-å·²ç»æ”¾å‡º',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `kkid` (`kkid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;



# Dump of table t_list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_list`;

CREATE TABLE `t_list` (
  `uid` int(14) unsigned NOT NULL COMMENT 'ç”¨æˆ·UID',
  `list_uuid` varchar(64) NOT NULL DEFAULT '',
  `fuid` bigint(20) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table t_report
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_report`;

CREATE TABLE `t_report` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `msg` varchar(255) NOT NULL DEFAULT '',
  `itime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `porcflag` int(4) NOT NULL DEFAULT '0',
  `retmsg` varchar(255) NOT NULL DEFAULT '',
  `reply` int(4) NOT NULL DEFAULT '0',
  `app_key` varchar(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_team_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_team_info`;

CREATE TABLE `t_team_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `creater` bigint(20) unsigned NOT NULL DEFAULT '0',
  `teamid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `name` varchar(64) NOT NULL DEFAULT '',
  `type` int(11) unsigned NOT NULL DEFAULT '0',
  `maxnum` int(11) unsigned NOT NULL DEFAULT '50',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `coreinfo` varchar(100) DEFAULT '',
  `exinfo` varchar(255) DEFAULT '',
  `del_flag` int(4) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `teamid_index` (`teamid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_team_list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_team_list`;

CREATE TABLE `t_team_list` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `teamid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `uid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `itime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `tid_uid` (`teamid`,`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_user_info
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_user_info`;

CREATE TABLE `t_user_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `did` varchar(64) NOT NULL DEFAULT '',
  `reg_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_date` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `baseinfo` varchar(255) DEFAULT '',
  `exinfo` varchar(255) NOT NULL DEFAULT '',
  `phonenum` varchar(16) NOT NULL DEFAULT '',
  `password` varchar(33) NOT NULL DEFAULT '',
  `platform` varchar(8) NOT NULL DEFAULT '',
  `setupid` bigint(20) unsigned NOT NULL DEFAULT '0',
  `v` bigint(20) unsigned NOT NULL DEFAULT '0',
  `bv` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_2` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table t_whitelist
# ------------------------------------------------------------

DROP TABLE IF EXISTS `t_whitelist`;

CREATE TABLE `t_whitelist` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `fuid` bigint(20) NOT NULL DEFAULT '0',
  `itime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_flag` int(1) NOT NULL DEFAULT '0' COMMENT '0 ï¼šæ­£å¸¸ï¼Œ1ï¼šå·²åˆ é™¤',
  `last_modify_date` int(14) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_fuid` (`uid`,`fuid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_user_login_info`;

CREATE TABLE `t_user_login_info` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` BIGINT(20) NOT NULL DEFAULT '0',
  `password` VARCHAR(33) NOT NULL DEFAULT '',
  `pc_login_core_token` VARCHAR(256) DEFAULT NULL,
  `pc_login_token_encrypted` VARCHAR(256) DEFAULT NULL,
  `pc_setup_id` VARCHAR(128) DEFAULT NULL,
  `pc_device_id` VARCHAR(128) DEFAULT NULL,
  `pc_time` BIGINT(20) NOT NULL DEFAULT 0,
  `mobile_login_core_token` VARCHAR(256) DEFAULT NULL,
  `mobile_login_token_encrypted` VARCHAR(256) DEFAULT NULL,
  `mobile_setup_id` VARCHAR(128) DEFAULT NULL,
  `mobile_device_id` VARCHAR(128) DEFAULT NULL,
  `mobile_platform_type` VARCHAR(33) NOT NULL DEFAULT '',
  `mobile_time` BIGINT(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`)
) ENGINE=MYISAM DEFAULT CHARSET=utf8;


/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
