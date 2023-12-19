insert into metric (type, name, value)
values ($1, $2, $3)
on conflict(type, name) do update set value = EXCLUDED.value
returning type, name, value;