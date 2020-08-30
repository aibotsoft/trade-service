create or
alter proc dbo.uspSaveEventPeriod @EventId varchar(100),
                                  @PeriodCode varchar(100),
                                  @IsActive bit
as
begin
    set nocount on
    MERGE dbo.EventPeriod AS t
    USING (select @EventId, @PeriodCode) s (EventId, PeriodCode)
    ON (t.EventId = s.EventId and t.PeriodCode = s.PeriodCode)

    WHEN MATCHED THEN
        UPDATE
        SET IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventId, PeriodCode, IsActive)
        VALUES (s.EventId, s.PeriodCode, @IsActive);
end