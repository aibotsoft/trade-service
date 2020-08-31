create table dbo.Surebet
(
    Id            int identity,
    HomeBetslipId varchar(32)                                   not null,
    AwayBetslipId varchar(32)                                   not null,
    Home          decimal(9, 5)                                 not null,
    Away          decimal(9, 5)                                 not null,
    Margin        decimal(9, 5)                                 not null,

    HomeReal      decimal(9, 5)                                 not null,
    AwayReal      decimal(9, 5)                                 not null,
    Profit        decimal(9, 5)                                 not null,
    HomeMax       decimal(9, 4)                                 not null,
    AwayMax       decimal(9, 4)                                 not null,

    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Surebet primary key (Id),
--     constraint UK_BetSlip unique (EventId, SportCode, BetType),
)