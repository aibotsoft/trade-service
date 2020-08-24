create table dbo.Auth
(
    Session varchar(50) not null,
    CreatedAt   datetimeoffset(2) default sysdatetimeoffset() not null,
    LastCheckAt datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_Auth primary key (CreatedAt),
)