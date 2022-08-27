CREATE TABLE IF NOT EXISTS Orders (
    Id varchar(64) PRIMARY KEY,
    AppUser varchar(32) DEFAULT NULL,
    AppTransId varchar(32) DEFAULT NULL,
    Item varchar(256) DEFAULT NULL,
    CreateAt bigint(20) DEFAULT NULL,
    TotalPrice bigint(20) DEFAULT NULL,
    Status varchar(32) DEFAULT NULL
);