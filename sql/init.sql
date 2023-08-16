create database if not exists tiktok
    charset utf8
    collate utf8_bin;

use tiktok;

create table if not exists user
(
    id               bigint        not null auto_increment,
    created_at       bigint        not null,
    deleted_at       bigint        not null default 0,

    username             varchar(32)   not null,
    password         varchar(32)   not null,
    avatar           varchar(1024) not null default '',
    background_image varchar(1024) not null default '',
    signature        varchar(1024) not null default '',
    primary key pk (id),
    unique key uk(username)
);

create table  if not exists video
(
    id             bigint        not null auto_increment,
    created_at     bigint        not null,
    updated_at     bigint        not null default 0,
    deleted_at     bigint        not null default 0,

    author_id      bigint        not null,
    title          varchar(32)   not null,
    play_url       varchar(1024) not null default '',
    cover_url      varchar(1024) not null default '',

    is_favorite    tinyint(1)    not null default 0,
    primary key pk (id)
);

create table if not exists  comment
(
    id         bigint        not null auto_increment,
    created_at bigint        not null,
    user_id    bigint        not null default 0,
    video_id   bigint        not null default 0,
    content    varchar(2048) not null,
    primary key pk (id)
);


create table  if not exists video_favor
(
    id         bigint not null auto_increment,
    created_at bigint not null,
    video_id   bigint not null,
    user_id    bigint not null,
    primary key pk (id),
    unique key uk (user_id, video_id)
);

create table  if not exists follow
(
    id         bigint not null auto_increment,
    created_at bigint not null,
    followee_id   bigint default 0 not null comment '被关注的人',
    follower_id bigint default 0 not null  comment '关注者',
    primary key pk(id),
    index ifee (follower_id),
    index ifer (followee_id)
);




