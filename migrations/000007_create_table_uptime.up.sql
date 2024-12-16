create table uptime (
	id serial not null,
	created_at timestamp with time zone not null default current_timestamp,
	primary key(id)
)
