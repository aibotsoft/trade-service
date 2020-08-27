create or alter proc dbo.uspSaveTeam @Id int,
                           @Name varchar(500)
as
begin
    set nocount on
    MERGE dbo.Team AS t
    USING (select @Id) s (Id)
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name  = @Name,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id, Name)
        VALUES (s.Id, @Name);
end