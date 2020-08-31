create or
alter proc dbo.uspSavePrice @BetslipId varchar(32),
                            @Bookie varchar(32),
                            @BetType varchar(32),
                            @Num tinyint,
                            @Price decimal(9, 5),
                            @Min decimal(9, 4),
                            @Max decimal(9, 4),
                            @Status varchar(32),
                            @IsActive bit
as
begin
    set nocount on
    MERGE dbo.Price AS t
    USING (select @BetslipId, @Bookie, @BetType, @Num) s (BetslipId, Bookie, BetType, Num)
    ON (t.BetslipId = s.BetslipId and t.Bookie = s.Bookie and t.BetType = s.BetType and t.Num = s.Num)

    WHEN MATCHED THEN
        UPDATE
        SET Status    = @Status,
            Price     = @Price,
            Min       = @Min,
            Max       = @Max,
            IsActive  = @IsActive,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (BetslipId, Bookie, BetType, Num, Price, Min, Max, Status, IsActive)
        VALUES (s.BetslipId, s.Bookie, s.BetType, s.Num,  @Price, @Min, @Max, @Status, @IsActive);
end