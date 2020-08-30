create or
alter proc dbo.uspSaveDoubleChance @EventPeriodId int,
                               @AwayDraw decimal(9, 5),
                               @HomeAway decimal(9, 5),
                               @HomeDraw decimal(9, 5),
                               @Margin decimal(9, 5),
                               @IsActive bit
as
begin
    set nocount on
    MERGE dbo.DoubleChance AS t
    USING (select @EventPeriodId) s (EventPeriodId)
    ON (t.EventPeriodId = s.EventPeriodId)

    WHEN MATCHED THEN
        UPDATE
        SET AwayDraw      = @AwayDraw,
            HomeAway      = @HomeAway,
            HomeDraw      = @HomeDraw,
            Margin    = @Margin,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, AwayDraw, HomeAway, HomeDraw, Margin, IsActive)
        VALUES (s.EventPeriodId, @AwayDraw, @HomeAway, @HomeDraw, @Margin, @IsActive);
end