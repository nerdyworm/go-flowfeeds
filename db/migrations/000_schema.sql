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
  user_id bigint references users(id) not null,
  episode_id bigint references episode(id) not null
);
create index index_listen_episode_id on listen using btree(episode_id);

create table favorite (
  id bigserial not null primary key,
  user_id bigint references users(id) not null,
  episode_id bigint references episode(id) not null
);
create index index_favorite_episode_id on favorite using btree(episode_id);
create unique index index_favorite_user_id_episode_id_unique on favorite using btree(user_id, episode_id);


CREATE OR REPLACE FUNCTION update_listens_count() RETURNS TRIGGER AS $update_listens_trigger$
BEGIN
  IF (TG_OP = 'DELETE') THEN
    UPDATE episode SET listens_count = (select count(*) from listen where episode_id = OLD.episode_id) WHERE episode.id = OLD.episode_id;
    RETURN OLD;
  ELSIF (TG_OP = 'INSERT') THEN
    UPDATE episode SET listens_count = (select count(*) from listen where episode_id = NEW.episode_id) WHERE episode.id = NEW.episode_id;
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
    UPDATE episode SET favorites_count = (select count(*) from favorite where episode_id = OLD.episode_id) WHERE episode.id = OLD.episode_id;
    RETURN OLD;
  ELSIF (TG_OP = 'INSERT') THEN
    UPDATE episode SET favorites_count = (select count(*) from favorite where episode_id = NEW.episode_id) WHERE episode.id = NEW.episode_id;
    RETURN NEW;
  END IF;
  RETURN NULL;
END;
$update_favorites_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER update_favorites_trigger
AFTER INSERT OR DELETE ON favorite FOR EACH ROW EXECUTE PROCEDURE update_favorites_count();
