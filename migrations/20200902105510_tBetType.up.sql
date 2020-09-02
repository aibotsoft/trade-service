create table dbo.BetType
(
    Id   tinyint     not null,
    Code varchar(32) not null,
    constraint PK_BetType primary key (Id),
)