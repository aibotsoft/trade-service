create table dbo.Duo
(
    EventPeriodId int                                           not null,
    BetTypeId     tinyint                                       not null,
    Code          smallint                                      not null,
    APrice        decimal(9, 5)                                 not null,
    BPrice        decimal(9, 5)                                 not null,
    IsActive      bit                                           not null default 1,
    CreatedAt     datetimeoffset(4) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(4) default sysdatetimeoffset() not null,
    constraint PK_Duo primary key (EventPeriodId, BetTypeId, Code),
)