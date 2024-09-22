-- 数据库
CREATE DATABASE kshared;

CREATE TABLE author (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,          -- 自增主键
    name VARCHAR(255) NOT NULL                           -- 作者名字，最大长度为255字符
);

CREATE TABLE article (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,          -- 自增主键
    title VARCHAR(255) NOT NULL,                         -- 标题，最大长度为255字符
    type VARCHAR(100) NOT NULL,                          -- 类型，最大长度为100字符
    label VARCHAR(100),                                  -- 标签，最大长度为100字符
    release_time DATETIME NOT NULL,                      -- 发布时间，使用 DATETIME 类型
    update_time DATETIME NOT NULL,                       -- 更新时间，使用 DATETIME 类型
    author_id INT UNSIGNED,                              -- 外键，引用 author 表的 id
    FOREIGN KEY (author_id) REFERENCES author(id)        -- 定义外键约束
);


INSERT INTO author (name) VALUES
('ty'),
('niuniu'),
('ysc');

INSERT INTO article (title, type, label, release_time, update_time, author_id) VALUES
('Go-Web net/http', 'modern pragram language', 'go net http multiplex router', '2024-09-01 10:00:00', '2024-09-02 12:00:00', 1),
('fetch-api', 'modern pragram language', 'js fetch async co-routine', '2024-08-20 09:30:00', '2024-08-25 14:45:00', 1),
('crud instance', 'modern pragram language', 'instance', '2024-07-15 08:15:00', '2024-07-20 16:30:00', 1);

select article.id, title, name, update_time
        from article 
        join author on author_id = author.id
        order by update_time desc
        limit 2;

CREATE TABLE comment (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pid INT UNSIGNED,
    article_id INT UNSIGNED,
    sender_id INT UNSIGNED,
    content VARCHAR(1024),
    release_time DATETIME,
    FOREIGN KEY (article_id) REFERENCES article(id),
    FOREIGN KEY (sender_id) REFERENCES author(id)
);
    -- use virtual root, so control in server.
    -- FOREIGN KEY (pid) REFERENCES comments(id),

SELECT 1 FROM comment where id = 1;

select * from author;
select * from comment;

select name form author 
join comment on author.id = comment.sender_id
where comment.id = 0;

select comment.id, name, comment.release_time, content 
from comment join author on sender_id = author.id
where pid = 0 and article_id = 7;


-- Tag 
CREATE TABLE tag (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    UNIQUE(name)
);

CREATE TABLE article_tag (
    article_id INT UNSIGNED,
    tag_id INT UNSIGNED,
    PRIMARY KEY(article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES article(id),
    FOREIGN KEY (tag_id) REFERENCES tag(id)
);

INSERT INTO tag (name) VALUES 
('go'),
('net'),
('multiplex'),
('router'),
('js'),
('fetch'),
('async'),
('co-routine'),
('instance'),
('concurrent'),
('goroutinue'),
('network-io'),
("linux-c"),
('epoll'),
('rb-tree'),
('algorithm'),
('cpp');

SELECT * FROM tag ORDER BY id;
INSERT INTO article_tag (article_id, tag_id) VALUES
(4, 1),
(4, 2),
(4, 3),
(4, 4),
(5, 5),
(5, 6),
(5, 7),
(5, 8),
(6, 9),
(7, 1),
(7, 10),
(7, 8),
(7, 11),
(7, 7),
(8, 12),
(8, 7),
(8, 13),
(8, 14),
(8, 15),
(9, 16),
(9, 17),
(9, 1);

SELECT name FROM article_tag 
JOIN article ON article_id = article.id 
JOIN tag ON tag_id = tag.id 
where article_id = 15;

SELECT name FROM tag WHERE name = ' ';

CREATE TABLE notify (
    id UINT32 PRIMARY KEY,
    source_id UINT32,
    target_id UINT32,
    comment_id UINT32,
    status ENUM('unread', 'read', 'deleted') NOT NULL DEFAULT 'unread'
);

CREATE TABLE notify (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source_id INT UNSIGNED,
    target_id INT UNSIGNED,
    comment_id INT UNSIGNED,
    status ENUM('unread', 'read', 'deleted') NOT NULL DEFAULT 'unread',
    FOREIGN KEY (source_id) REFERENCES author(id),
    FOREIGN KEY (target_id) REFERENCES author(id),
    FOREIGN KEY (comment_id) REFERENCES comment(id)
);