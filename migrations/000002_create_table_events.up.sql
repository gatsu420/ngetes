create table events (
	id serial not null,
	task_id integer not null,
	name text not null,
	created_at timestamp with time zone not null default current_timestamp,
	primary key (id)
)
