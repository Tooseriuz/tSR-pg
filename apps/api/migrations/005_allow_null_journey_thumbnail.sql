alter table journeys
  alter column thumbnail drop not null;

---- create above / drop below ----

alter table journeys
  alter column thumbnail set not null;
