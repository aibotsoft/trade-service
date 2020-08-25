create table dbo.Account
(
    Id           int identity  not null,
--     AccountType  varchar(100)  not null default 'BET',
--     CurrencyCode varchar(100)  not null default 'USD',
--     ServiceName  varchar(100)  not null,
    Username     varchar(100)  not null,
    Password     varchar(100)  not null,
--     Commission   decimal(9, 5) not null default 0,
--     Share        decimal(9, 5) not null default 1,

    CreatedAt    datetimeoffset(2)         default sysdatetimeoffset() not null,
    UpdatedAt    datetimeoffset(2)         default sysdatetimeoffset() not null,

    constraint PK_Account primary key (Id),
--     constraint CH_AccountType check (AccountType in ('SCAN', 'BET')),
--     constraint CH_CurrencyCode check (CurrencyCode in ('USD', 'EUR')),
)