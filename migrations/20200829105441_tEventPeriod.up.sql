create table dbo.EventPeriod
(
    Id         int identity                                  not null,
    EventId    varchar(100)                                  not null,
    PeriodCode varchar(100)                                  not null,
    CreatedAt  datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt  datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_EventPeriod primary key (Id),
)