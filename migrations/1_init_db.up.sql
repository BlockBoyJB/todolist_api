create table if not exists public.user
(
    id         bigserial primary key,
    username   varchar unique not null,
    password   varchar        not null,
    created_at timestamp      not null default now()
);

create table if not exists task
(
    id          bigserial primary key,
    username    varchar references public.user (username),
    title       varchar   not null,
    description text      not null,
    due_date    timestamp not null,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now()
);

create or replace function update_time()
    returns trigger as
$$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;

create trigger update_time_trigger
    before update
    on task
    for each row
execute function update_time();
