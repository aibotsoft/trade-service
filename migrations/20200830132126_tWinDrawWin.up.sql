create table dbo.WinDrawWin
(
    EventPeriodId int                                           not null,
    Away          decimal(9, 5)                                 not null,
    Home          decimal(9, 5)                                 not null,
    Draw          decimal(9, 5)                                 not null,
    Margin        decimal(9, 5)                                 not null,
    IsActive      bit                                           not null default 1,
    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_WinDrawWin primary key (EventPeriodId),
)