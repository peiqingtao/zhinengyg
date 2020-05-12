


create database projectA charset utf8mb4;
use projectA;

create user projectAUser@localhost identified by 'hellokang';
grant all privileges on projectA.* to projectAUser@localhost;
flush privileges;

create table if not exists a_categories (
  id int unsigned auto_increment,
  parent_id int unsigned,
  name varchar(255),
  logo varchar(255),
  description varchar(255),
  sort_order int,
  meta_title varchar(255),
  meta_keywords varchar(255),
  meta_description varchar(255),
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp,
  primary key (id),
  index (parent_id),
  index (name),
  index (sort_order),
  index (updated_at),
  index (deleted_at)
)engine innodb charset utf8mb4;
alter table a_categories add column created_at timestamp, add column updated_at timestamp, add column deleted_at timestamp;

insert into a_categories (id, name, parent_id) values (1, '未分类', 0);
insert into a_categories (id, name, parent_id) values (2, '图书', 0);
insert into a_categories (id, name, parent_id) values (3, '电脑', 0);
insert into a_categories (id, name, parent_id) values (4, '纸质书', 2);
insert into a_categories (id, name, parent_id) values (5, '电子书', 2);
insert into a_categories (id, name, parent_id) values (6, '笔记本', 3);
insert into a_categories (id, name, parent_id) values (7, '平板', 3);
insert into a_categories (id, name, parent_id) values (8, '一体机', 3);
insert into a_categories (id, name, parent_id) values (9, '13英寸', 6);
insert into a_categories (id, name, parent_id) values (10, '14英寸', 6);
insert into a_categories (id, name, parent_id) values (11, '15英寸', 6);
insert into a_categories (id, name, parent_id) values (12, '17英寸', 6);