create table dbo.Price
(
    BetslipId varchar(32)                                   not null,
    Bookie    varchar(32)                                   not null,
    BetType   varchar(32)                                   not null,
    Status    varchar(32)                                   not null,
    Num       tinyint                                       not null,
    Price     decimal(9, 5)                                 not null,
    Min       decimal(9, 5)                                 not null,
    Max       decimal(9, 5)                                 not null,
    IsActive  bit                                           not null default 1,
    CreatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Price primary key (BetslipId, Bookie, BetType),
)