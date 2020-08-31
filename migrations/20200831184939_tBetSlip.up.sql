create table dbo.BetSlip
(
    Id               int identity,
    BetslipId        varchar(32)                                   not null,
    EventId          varchar(100)                                  not null,
    SportCode        varchar(100)                                  not null,
    BetType          varchar(32)                                   not null,

    BetTypeDes       varchar(1000)                                 not null,
    BetTypeTemp      varchar(1000)                                 not null,

    EquivalentBets   bit                                           not null,
    MultipleAccounts bit                                           not null,
    IsOpen           bit                                           not null,
    ExpiryTs         float                                         not null,

    CreatedAt        datetimeoffset(2) default sysdatetimeoffset() not null,
    UpdatedAt        datetimeoffset(2) default sysdatetimeoffset() not null,
    constraint PK_BetSlip primary key (Id),
    constraint UK_BetSlip unique (EventId, SportCode, BetType),
)