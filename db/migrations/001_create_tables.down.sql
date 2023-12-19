begin transaction;

-- Удаляем тип метрики
drop type if exists metric_type;

-- Удаляем метрику
drop table if exists metric;

commit;