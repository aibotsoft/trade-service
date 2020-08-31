create table dbo.Surebet
(
    SurebetId     bigint                                        not null,
    BetslipId     varchar(32)                                   not null,
    Price         decimal(9, 5)                                 not null,

    BestPrice     decimal(9, 5)                                 not null,
    WeightedPrice decimal(9, 5)                                 not null,
    Min           decimal(9, 4)                                 not null,
    Max           decimal(9, 4)                                 not null,
    Volume        decimal(9, 2)                                 not null,
    Bookie        varchar(32)                                   not null,
    PriceCount    tinyint                                       not null,

    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Surebet primary key (SurebetId, BetslipId),
)