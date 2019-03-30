CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists users (
  id uuid default uuid_generate_v4(),
  login varchar(256) unique not null,
  password varchar(512) not null,
  created timestamp default NOW(),

  primary key(id)
);

create table if not exists links (
  id uuid default uuid_generate_v4(),
  url text not null,
  created timestamp default NOW(),
  user_id uuid not null,

  primary key(id),
  constraint links_user_id foreign key (user_id) references users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create table if not exists usages (
  id uuid default uuid_generate_v4(),
  created timestamp default NOW(),
  meta json default null,
  link_id uuid not null,

  primary key(id),
  constraint usages_link_id foreign key (link_id) references links(id) ON DELETE CASCADE ON UPDATE CASCADE
);

insert into users (login, password) values('anon', '');
