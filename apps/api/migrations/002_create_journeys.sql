create table journeys (
  id bigserial primary key,
  name text not null,
  timestamp date not null,
  location varchar(30) not null,
  thumbnail text not null,
  created_at date not null default current_date
);

---- create above / drop below ----

drop table journeys;
