create table journey_images (
  id uuid primary key default gen_random_uuid(),
  path text not null,
  created_at date not null default current_date
);

---- create above / drop below ----

drop table journey_images;
