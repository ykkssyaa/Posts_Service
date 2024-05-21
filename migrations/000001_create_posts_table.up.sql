CREATE TABLE IF NOT EXISTS Posts
(
    id serial primary key,
    created_at timestamp default now(),
    name varchar(100),
    content varchar(2000),
    author varchar(100),
    comments_allowed boolean
);