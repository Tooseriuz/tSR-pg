alter table journeys
  add column highlight boolean not null default false;

---- create above / drop below ----

alter table journeys
  drop column highlight;
