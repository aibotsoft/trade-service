create table dbo.League
(
    Id        int                                        not null,
    Name      varchar(1000)                              not null,
    Country   varchar(300),
    Sport     varchar(300),
    CreatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_League primary key (Id),
)