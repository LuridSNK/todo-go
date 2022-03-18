/* todo: try to utilize https://www.unixtimestamp.com/ for naming migrations ? */

create table if not exists TodoItems (
    id uuid primary key default gen_random_uuid(),
    description text,
    createdAt timestamp without time zone default (now() at time zone 'utc'),
	isDone bool);

---- create above / drop below ----
    
drop table TodoItems;