create or
alter proc dbo.uspSaveWinDrawWin @EventPeriodId int,
                                   @Away decimal(9, 5),
                                   @Home decimal(9, 5),
                                   @Draw decimal(9, 5),
                                   @Margin decimal(9, 5),
                                   @IsActive bit
as
begin
    set nocount on
    MERGE dbo.WinDrawWin AS t
    USING (select @EventPeriodId) s (EventPeriodId)
    ON (t.EventPeriodId = s.EventPeriodId)

    WHEN MATCHED THEN
        UPDATE
        SET Away      = @Away,
            Home      = @Home,
            Draw      = @Draw,
            Margin    = @Margin,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, Away, Home, Draw, Margin, IsActive)
        VALUES (s.EventPeriodId, @Away, @Home, @Draw, @Margin, @IsActive);
end