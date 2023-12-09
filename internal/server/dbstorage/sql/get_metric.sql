select type, name, delta, value
from metric
where type = $1
    and name = $2;