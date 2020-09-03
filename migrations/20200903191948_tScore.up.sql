create table dbo.ScoreFootball
(
    EventPeriodId int not null,
    RedHome       tinyint,
    RedAway       tinyint,
    ScoreHome     tinyint,
    ScoreAway     tinyint,
    PeriodCode    varchar(32),
    PeriodMin     tinyint,
    constraint PK_ScoreFootball primary key (EventPeriodId),
)