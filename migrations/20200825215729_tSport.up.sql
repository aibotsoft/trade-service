create table dbo.Sport
(
    Id        smallint                                      not null,
    Name      varchar(1000)                                 not null,
    CreatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Sport primary key (Id),
)