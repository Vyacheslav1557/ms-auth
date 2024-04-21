create table if not exists users
(
    id              bigint primary key generated always as identity,
    username        text unique,
    hashed_password text,
    is_logged_in    bool,
    last_login_at   bigint
);