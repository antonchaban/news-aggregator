create table sources
(
    id   bigserial
        primary key,
    name varchar(2048),
    link varchar(16384)
);

create table articles
(
    id          bigserial
        primary key,
    title       varchar(2048),
    description varchar(16384),
    link        varchar(16384),
    source_id   bigint
        references sources (id) on delete cascade,
    pub_date    timestamp default now()
);