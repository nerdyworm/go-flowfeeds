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
  feed bigint references feed(id) not null,
  guid text not null,
  title text not null,
  description text not null,
  url text not null,
  image text not null,
  published timestamp with time zone not null,
  listens_count int not null default 0,
  favorites_count int not null default 0
);

create unique index episodes_guid_unique on episode using btree(guid);

create table users (
  id bigserial not null primary key,
  email text not null,
  encrypted_password text not null
);

create unique index user_email_unique on users using btree(email);

create table listen (
  id bigserial not null primary key,
  "user" bigint references users(id) not null,
  episode bigint references episode(id) not null
);
create index index_listen_episode on listen using btree(episode);

create table favorite (
  id bigserial not null primary key,
  "user" bigint references users(id) not null,
  episode bigint references episode(id) not null
);
create index index_favorite_episode on favorite using btree(episode);
create unique index index_favorite_user_episode_unique on favorite using btree("user", episode);


CREATE OR REPLACE FUNCTION update_listens_count() RETURNS TRIGGER AS $update_listens_trigger$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    UPDATE episode SET listens_count = (select count(*) from listen where episode = OLD.episode) WHERE episode.id = OLD.episode;
    RETURN OLD;
  ELSIF (TG_OP = 'INSERT') THEN
    UPDATE episode SET listens_count = (select count(*) from listen where episode = NEW.episode) WHERE episode.id = NEW.episode;
    RETURN NEW;
  END IF;
  RETURN NULL;
END;
$update_listens_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER update_listens_trigger
AFTER INSERT OR DELETE ON listen FOR EACH ROW EXECUTE PROCEDURE update_listens_count();

CREATE OR REPLACE FUNCTION update_favorites_count() RETURNS TRIGGER AS $update_favorites_trigger$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    UPDATE episode SET favorites_count = (select count(*) from favorite where episode = OLD.episode) WHERE episode.id = OLD.episode;
    RETURN OLD;
  ELSIF (TG_OP = 'INSERT') THEN
    UPDATE episode SET favorites_count = (select count(*) from favorite where episode = NEW.episode) WHERE episode.id = NEW.episode;
    RETURN NEW;
  END IF;
  RETURN NULL;
END;
$update_favorites_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER update_favorites_trigger
AFTER INSERT OR DELETE ON favorite FOR EACH ROW EXECUTE PROCEDURE update_favorites_count();
