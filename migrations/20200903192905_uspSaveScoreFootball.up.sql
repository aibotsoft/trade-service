create or
alter proc dbo.uspSaveScoreFootball @EventPeriodId int,
                                    @RedHome tinyint = null,
                                    @RedAway tinyint = null,
                                    @ScoreHome tinyint = null,
                                    @ScoreAway tinyint = null,
                                    @PeriodCode varchar(32) = null,
                                    @PeriodMin tinyint = null
as
begin
    set nocount on
    MERGE dbo.ScoreFootball AS t
    USING (select @EventPeriodId) s (EventPeriodId)
    ON (t.EventPeriodId = s.EventPeriodId)

    WHEN MATCHED THEN
        UPDATE
        SET RedHome    = @RedHome,
            RedAway    = @RedAway,
            ScoreHome  = @ScoreHome,
            ScoreAway  = @ScoreAway,
            PeriodCode = @PeriodCode,
            PeriodMin  = @PeriodMin,
            UpdatedAt  = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (EventPeriodId, RedHome, RedAway, ScoreHome, ScoreAway, PeriodCode, PeriodMin)
        VALUES (s.EventPeriodId, @RedHome, @RedAway, @ScoreHome, @ScoreAway, @PeriodCode, @PeriodMin);
end