create or
alter proc dbo.uspSaveTrio @EventPeriodId int,
                           @BetTypeId tinyint,
                           @APrice decimal(9, 5),
                           @BPrice decimal(9, 5),
                           @CPrice decimal(9, 5),
                           @IsActive bit
as
begin
    set nocount on
    MERGE dbo.Trio AS t
    USING (select @EventPeriodId, @BetTypeId) s (EventPeriodId, BetTypeId)
    ON (t.EventPeriodId = s.EventPeriodId and t.BetTypeId = s.BetTypeId)

    WHEN MATCHED THEN
        UPDATE
        SET APrice    = @APrice,
            BPrice    = @BPrice,
            CPrice    = @CPrice,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, BetTypeId, APrice, BPrice, CPrice, IsActive)
        VALUES (s.EventPeriodId, @BetTypeId, @APrice, @BPrice, @CPrice, @IsActive);
end