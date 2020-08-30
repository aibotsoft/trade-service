create table dbo.Event
(
    Id        varchar(100)                                  not null,
    HomeId    int                                           not null,
    AwayId    int                                           not null,
    LeagueId  int                                           not null,
    Starts    datetimeoffset(2)                             not null,
    CreatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Event primary key (Id),
)