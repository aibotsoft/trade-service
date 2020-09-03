create table dbo.ScoreFootball
(
    EventPeriodId int not null,
    RedHome       tinyint,
    RedAway       tinyint,
    ScoreHome     tinyint,
    ScoreAway     tinyint,
    PeriodCode    varchar(32),
    PeriodMin     tinyint,
    CreatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt     datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_ScoreFootball primary key (EventPeriodId),
)