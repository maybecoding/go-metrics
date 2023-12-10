begin transaction;

-- Тип метрики
drop type if exists metric_type;
create type metric_type as enum('gauge', 'counter');

-- Метрика
drop table if exists metric;
create table metric (
    id int primary key generated always as identity,
    type metric_type not null,
    name varchar(255) not null,
    value double precision null,
    delta int8 null
);

create unique index IX_metric_type_name on metric (type, name);

alter table metric
    add constraint CH_metric__gauge_only_value__counter_only_delta check(
                    type = 'gauge' and value is not null and delta is null or
                    type = 'counter' and value is null and delta is not null
        );

commit;