CREATE TABLE IF NOT EXISTS Comments
(
    id serial primary key,
    created_at timestamp default now(),
    content varchar(2000),
    author varchar(100),
    post int not null,
    FOREIGN KEY (post) REFERENCES Posts(id) ON DELETE CASCADE ON UPDATE CASCADE ,
    reply_to int,
    FOREIGN KEY (reply_to) REFERENCES Comments(id) ON DELETE SET NULL ON UPDATE CASCADE
);