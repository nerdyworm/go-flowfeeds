create table feed (
  id bigserial not null primary key,
  url text not null,
  title text not null,
  description text not null,
  image text not null,
  updated timestamp with time zone not null
);
create unique index feeds_url_unique on feed using btree(url);

create table episode (
  id bigserial not null primary key,
  feed_id bigint references feed(id) not null,
  guid text not null,
  title text not null,
  description text not null,
  url text not null,
  image text not null,
  published timestamp with time zone not null
);

create unique index episodes_guid_unique on episode using btree(guid);
