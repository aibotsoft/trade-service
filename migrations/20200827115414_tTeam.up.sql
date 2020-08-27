create table dbo.Team
(
    Id        int                                           not null,
    Name      varchar(500),
    CreatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Team primary key (Id),
);