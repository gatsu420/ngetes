create table tasks (
	id serial not null,
	name text not null,
	status text not null default 'backlog',
	created_at timestamp with time zone not null default current_timestamp,
	updated_at timestamp with time zone not null default current_timestamp,
	primary key (id)
)
