create or
alter proc dbo.uspSaveDoubleChance @EventPeriodId int,
                               @HandicapCode smallint,
                               @AwayDraw decimal(9, 5),
                               @HomeAway decimal(9, 5),
                               @HomeDraw decimal(9, 5),
                               @Margin decimal(9, 5),
                               @IsActive bit
as
begin
    set nocount on
    MERGE dbo.DoubleChance AS t
    USING (select @EventPeriodId, @HandicapCode) s (EventPeriodId, HandicapCode)
    ON (t.EventPeriodId = s.EventPeriodId and t.HandicapCode = s.HandicapCode)

    WHEN MATCHED THEN
        UPDATE
        SET AwayDraw      = @AwayDraw,
            HomeAway      = @HomeAway,
            HomeDraw      = @HomeDraw,
            Margin    = @Margin,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, HandicapCode, AwayDraw, HomeAway, HomeDraw, Margin, IsActive)
        VALUES (s.EventPeriodId, s.HandicapCode, @AwayDraw, @HomeAway, @HomeDraw, @Margin, @IsActive);
end