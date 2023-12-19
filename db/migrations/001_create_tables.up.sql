begin transaction;

-- Тип метрики
drop type if exists metric_type;
create type metric_type as enum('gauge', 'counter');

-- Метрика
drop table if exists metric;
create table metric (
    type metric_type not null,
    name varchar(255) not null,
    value double precision null,
    delta int8 null,
    constraint PK_metric_type_name primary key (type, name)
);

alter table metric
    add constraint CH_metric__gauge_only_value__counter_only_delta check(
                    type = 'gauge' and value is not null and delta is null or
                    type = 'counter' and value is null and delta is not null
        );

commit;