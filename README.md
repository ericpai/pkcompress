# pkcompress

[![Build Status](https://travis-ci.org/ericpai/pkcompress.svg?branch=master)](https://travis-ci.org/ericpai/pkcompress) [![codecov](https://codecov.io/gh/ericpai/pkcompress/branch/master/graph/badge.svg)](https://codecov.io/gh/ericpai/pkcompress) [![MIT license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fericpai%2Fpkcompress.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fericpai%2Fpkcompress?ref=badge_shield)

PkCompress is a tool to compress discrete integer primary keys of MySQL tables.


## Scenario

Assuming in your database there are three tables with the following definitions:

```mysql
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
```

And with those data:

```mysql
INSERT INTO `region`(`id`, `name`) VALUES (1000, 'East');
INSERT INTO `region`(`id`, `name`) VALUES (2000, 'West');

INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (100, 'Sales', 1000);
INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (200, 'Tech', 2000);

INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (10, 'Mr. Little', 100);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (20, 'Mrs. Medium', 100);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (30, 'Mr. Big', 200);
```

Now we want to make all the primary keys of the rows to be continuous from 1 and the `AUTO_INCREMENT` to be the exact value. What we should do now? May we have the following choices:

- Update primary key directly from the first row: Note that we have foreign key constraints and the action will be prohibited as the update rule is `RESTRICT`.
- Update foreign keys before doing primary key: It's hard to write scripts to find all the foreign keys and the corresponding values between foreign and primary keys.
- Dump and re-insert the data: The reference relationship may be broken.

The PkCompress tool is used to solve this issue!

## Workflow

The PKCompress tool handles all the tables in one single database. As we think that foreign keys referenced to tables in other different database merely happens. And all the primary keys are 'logical' enough as they are single column with integer type and `AUTO_INCREMENT` property.

The PkCompress will work as the following procedures:

1. Connect to the database with the login parameters from the command line.
2. Get all the table names from `SHOW TABLES`.
3. Get all the primary keys of each table from table `information_schema.KEY_COLUMN_USAGE`.
4. Get the column type of each primary key from table `information_schema.columns`.
5. Get all the foreign keys and their update rules from table `information_schema.table_constraints` and `information_schema.referential_constraints` and `information_schema.KEY_COLUMN_USAGE`.
6. Alter all the foreign keys with update rule `CASCADE` with `ALTER` statement.
7. For each table, select all the primary keys ascent first and then update each primary key to be the natural number in a single transaction.
8. Alter the `AUTO_INCREMENT` property to the 'to be next' primary key.
9. Resume all the foreign keys' update rule with `ALTER` statement.

After handling successfully, we will see the data in the previous example as the following:

```mysql
INSERT INTO `region`(`id`, `name`) VALUES (1, 'East');
INSERT INTO `region`(`id`, `name`) VALUES (2, 'West');

INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (1, 'Sales', 1);
INSERT INTO `department`(`id`, `name`, `region_id`) VALUES (2, 'Tech', 2);

INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (1, 'Mr. Little', 1);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (2, 'Mrs. Medium', 1);
INSERT INTO `staff`(`id`, `name`, `department_id`) VALUES (3, 'Mr. Big', 2);
```

And the `AUTO_INCREMENT` value will be `3`, `3` and `4`.

## Installation

### From source code

You must have Go1.9 environment and [glide](https://github.com/Masterminds/glide) package management tool installed.

Run the following commands in your source directory.

```
glide install
go build -o pkcompress main.go
```

## Usage

As DDL statements can't be in transaction in MySQL, the constraints definitions may be changed if connecting error occurs during execution. So

**Please backup your database first!**

```
./pkcompress -h <mysql_host> -P <mysql_port> -u <mysql_user> -D <database_name>
```

Note that the user must have at least `SELECT, UPDATE, DROP, ALTER` privileges of all elements of this database.

After pressing `enter`, PkCompress will ask you to input the password interactively and secretly.

If it connects the database successfully, you can have a cup of coffee and watch the output while it is executing.


## Limitation

As the [MySQL document](https://dev.mysql.com/doc/refman/5.6/en/innodb-foreign-key-constraints.html) said, the update role will be treated as `RESTRICT` no matter what the constraint definition is if the referenced column is in the same table. So PkCompress can't handle this scenario yet.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fericpai%2Fpkcompress.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fericpai%2Fpkcompress?ref=badge_large)