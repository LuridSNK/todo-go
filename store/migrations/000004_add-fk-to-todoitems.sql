
alter table TodoItems add constraint fk_items_to_users foreign key (creatorId) references Users(Id);

---- create above / drop below ----

alter table TodoItems drop constraint fk_items_to_users;