create table sources
(
    id   bigserial
        primary key,
    name varchar(2048),
    link varchar(16384),
    short_name varchar(512) unique
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

insert into sources (name, link, short_name)
values ('BBC News', 'https://feeds.bbci.co.uk/news/rss.xml', 'bbc'),
       ('ABC News: International', 'https://abcnews.go.com/abcnews/internationalheadlines', 'abcnews'),
       ('The Washington Times stories: World', 'https://www.washingtontimes.com/rss/headlines/news/world/', 'washingtontimes'),
       ('USA TODAY', 'https://www.usatoday.com/news/world/', 'usatoday'),
       ('The Washington Post', 'http://feeds.washingtonpost.com/rss/world', 'washingtonpost');