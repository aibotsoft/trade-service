create table dbo.Total
(
    EventPeriodId int                                           not null,
    HandicapCode  smallint                                      not null,
    Over          decimal(9, 5)                                 not null,
    Under         decimal(9, 5)                                 not null,
    Margin        decimal(9, 5)                                 not null,
    IsActive      bit                                           not null default 1,
    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Total primary key (EventPeriodId, HandicapCode),
)