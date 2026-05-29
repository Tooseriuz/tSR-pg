create table journey_contents (
  journey_id bigint primary key references journeys(id) on delete cascade,
  content text not null
);

---- create above / drop below ----

drop table journey_contents;
