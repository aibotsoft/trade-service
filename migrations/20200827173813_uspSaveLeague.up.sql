create or alter proc dbo.uspSaveLeague @Id int,
                                    @Name varchar(1000),
                                    @Country varchar(300),
                                    @SportId smallint
as
begin
    set nocount on
    MERGE dbo.League AS t
    USING (select @Id) s (Id)
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name      = @Name,
            Country   = @Country,
            SportId   = @SportId,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id, Name, Country, SportId)
        VALUES (s.Id, @Name, @Country, @SportId);
end