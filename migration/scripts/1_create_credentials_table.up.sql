create table if not exists teacher (
    id uuid primary key,
    username varchar not null unique,
    password bytea not null,

    check(length(username) >= 4 and length(username) <= 30)
);
