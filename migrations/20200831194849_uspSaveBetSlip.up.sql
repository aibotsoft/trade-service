create or
alter proc dbo.uspSaveBetSlip @BetslipId varchar(32),
                              @EventId varchar(100),
                              @SportCode varchar(100),
                              @BetType varchar(32),
                              @BetTypeDes varchar(1000),
                              @BetTypeTemp varchar(1000),
                              @EquivalentBets bit,
                              @MultipleAccounts bit,
                              @IsOpen bit,
                              @ExpiryTs float
as
begin
    set nocount on
    MERGE dbo.BetSlip AS t
    USING (select @BetslipId) s (BetslipId)
    ON (t.BetslipId = s.BetslipId)

    WHEN MATCHED THEN
        UPDATE
        SET EventId    = @EventId,
            SportCode     = @SportCode,
            BetType       = @BetType,
            BetTypeDes       = @BetTypeDes,
            BetTypeTemp  = @BetTypeTemp,
            EquivalentBets  = @EquivalentBets,
            MultipleAccounts  = @MultipleAccounts,
            IsOpen  = @IsOpen,
            ExpiryTs  = @ExpiryTs,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (BetslipId, EventId, SportCode, BetType, BetTypeDes, BetTypeTemp, EquivalentBets, MultipleAccounts, IsOpen, ExpiryTs)
        VALUES (s.BetslipId, @EventId, @SportCode, @BetType, @BetTypeDes, @BetTypeTemp, @EquivalentBets, @MultipleAccounts, @IsOpen, @ExpiryTs);
end