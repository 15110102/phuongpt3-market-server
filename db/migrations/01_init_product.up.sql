CREATE TABLE IF NOT EXISTS Product (
    Id varchar(26) PRIMARY KEY,
    Name varchar(32) DEFAULT NULL,
    `Desc` varchar(32) DEFAULT NULL,
    CreateAt bigint(20) DEFAULT NULL,
    UpdateAt bigint(20) DEFAULT NULL,
    Price bigint(20) DEFAULT NULL,
    Image varchar(128) DEFAULT NULL
);