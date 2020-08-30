create or
alter proc dbo.uspSaveHandicap @EventPeriodId int,
                               @HandicapCode smallint,
                               @Away decimal(9, 5),
                               @Home decimal(9, 5),
                               @Margin decimal(9, 5),
                               @IsActive bit
as
begin
    set nocount on
    MERGE dbo.Handicap AS t
    USING (select @EventPeriodId, @HandicapCode) s (EventPeriodId, HandicapCode)
    ON (t.EventPeriodId = s.EventPeriodId and t.HandicapCode = s.HandicapCode)

    WHEN MATCHED THEN
        UPDATE
        SET Away      = @Away,
            Home      = @Home,
            Margin    = @Margin,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, HandicapCode, Away, Home, Margin, IsActive)
        VALUES (s.EventPeriodId, s.HandicapCode, @Away, @Home, @Margin, @IsActive);
end