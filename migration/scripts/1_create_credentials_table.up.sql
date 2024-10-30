create type account_type as enum ('TEACHER', 'STUDENT');

create table if not exists credentials (
    id uuid primary key,
    username varchar not null unique,
    password_hash bytea not null,
    account_type account_type not null,

    check(length(username) >= 4 and length(username) <= 30)
);
