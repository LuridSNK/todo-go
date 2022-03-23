
create table if not exists Users (
    id uuid primary key default gen_random_uuid(),
    email text,
    passwordHash text,
    createdAt timestamp without time zone default (now() at time zone 'utc'),
	isDeleted bool default false);

---- create above / drop below ----
    
drop table Users;