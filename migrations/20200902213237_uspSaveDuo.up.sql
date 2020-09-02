create or
alter proc dbo.uspSaveTrio @EventPeriodId int,
                           @BetTypeId tinyint,
                           @Code smallint,
                           @APrice decimal(9, 5),
                           @BPrice decimal(9, 5),
                           @IsActive bit
as
begin
    set nocount on
    MERGE dbo.Duo AS t
    USING (select @EventPeriodId, @BetTypeId, @Code) s (EventPeriodId, BetTypeId, Code)
    ON (t.EventPeriodId = s.EventPeriodId and t.BetTypeId = s.BetTypeId and t.Code = s.Code)

    WHEN MATCHED THEN
        UPDATE
        SET APrice    = @APrice,
            BPrice    = @BPrice,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, BetTypeId, Code, APrice, BPrice, IsActive)
        VALUES (s.EventPeriodId, s.BetTypeId, s.Code, @APrice, @BPrice, @IsActive);
end