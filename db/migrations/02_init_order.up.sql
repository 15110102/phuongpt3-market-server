CREATE TABLE IF NOT EXISTS Orders (
    Id int NOT NULL AUTO_INCREMENT,
    AppUser varchar(32) DEFAULT NULL,
    AppTransId varchar(64) DEFAULT NULL,
    ZpTransToken varchar(32) DEFAULT NULL,
    Item varchar(256) DEFAULT NULL,
    CreateAt bigint(20) DEFAULT NULL,
    TotalPrice bigint(20) DEFAULT NULL,
    Status varchar(32) DEFAULT NULL,
    PRIMARY KEY (Id)
);