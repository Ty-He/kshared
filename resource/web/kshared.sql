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
