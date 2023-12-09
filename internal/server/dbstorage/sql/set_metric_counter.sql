insert into metric (type, name, delta)
values ($1, $2, $3)
on conflict(type, name) do update set delta = EXCLUDED.delta
returning type, name, delta;