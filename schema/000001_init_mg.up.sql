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
    link        varchar(16384) unique,
    source_id   bigint
        references sources (id) on delete cascade,
    pub_date    timestamp default now()
);

insert into sources (name, link)
values ('BBC News', 'https://feeds.bbci.co.uk/news/rss.xml'),
       ('ABC News: International', 'https://abcnews.go.com/abcnews/internationalheadlines'),
       ('The Washington Times stories: World', 'https://www.washingtontimes.com/rss/headlines/news/world/'),
       ('USA TODAY', 'https://www.usatoday.com/news/world/'),
       ('The Washington Post', 'http://feeds.washingtonpost.com/rss/world');