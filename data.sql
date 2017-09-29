DROP DATABASE IF EXISTS `pktest`;
CREATE DATABASE `pktest`;
CREATE USER `pktest`@'%' IDENTIFIED by 'dsA12(djks';
GRANT ALL ON `pktest`.* TO `pktest`@'%';
FLUSH PRIVILEGES;

USE `pktest`;

CREATE TABLE `staff` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `department_id` int(11)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `department` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `region_id` int(11)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `region` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `staff` ADD CONSTRAINT fk_staff_department_id FOREIGN KEY (`department_id`) REFERENCES department(id) on UPDATE RESTRICT;
ALTER TABLE `department` ADD CONSTRAINT fk_department_region_id FOREIGN KEY (`region_id`) REFERENCES region(id) on UPDATE RESTRICT;

INSERT INTO `region`(`id`, `name`) VALUES (1000, 'East');
INSERT INTO `region`(`id`, `name`) VALUES (2000, 'West');

INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (100, 'Sales', 1000);
INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (200, 'Tech', 2000);

INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (10, 'Mr. Little', 100);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (20, 'Mrs. Medium', 100);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (30, 'Mr. Big', 200);