
alter table TodoItems add column creatorId uuid not null;

---- create above / drop below ----

alter table TodoItems drop column creatorId;