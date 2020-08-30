create table dbo.DoubleChance
(
    EventPeriodId int                                           not null,
    AwayDraw      decimal(9, 5)                                 not null,
    HomeAway      decimal(9, 5)                                 not null,
    HomeDraw      decimal(9, 5)                                 not null,
    Margin        decimal(9, 5)                                 not null,
    IsActive      bit                                           not null default 1,
    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_DoubleChance primary key (EventPeriodId),
)