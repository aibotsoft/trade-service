create or
alter proc dbo.uspSaveTotal @EventPeriodId int,
                               @HandicapCode smallint,
                               @Ove decimal(9, 5),
                               @Under decimal(9, 5),
                               @Margin decimal(9, 5),
                               @IsActive bit
as
begin
    set nocount on
    MERGE dbo.Total AS t
    USING (select @EventPeriodId, @HandicapCode) s (EventPeriodId, HandicapCode)
    ON (t.EventPeriodId = s.EventPeriodId and t.HandicapCode = s.HandicapCode)

    WHEN MATCHED THEN
        UPDATE
        SET Ove      = @Ove,
            Under      = @Under,
            Margin    = @Margin,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, HandicapCode, Ove, Under, Margin, IsActive)
        VALUES (s.EventPeriodId, s.HandicapCode, @Ove, @Under, @Margin, @IsActive);
end