create
or replace function calculate_distance(
    lat1 float,
    long1 float,
    lat2 float,
    long2 float
) returns float as 
$$ 
declare 
    earth_radius constant float := 6371.009;
    dlat float;
    dlong float;
    a float;
    c float;
    distance float;
begin 
dlat := (lat2 - lat1) * PI() / 180.0;

dlong := (long2 - long1) * PI() / 180.0;

lat1 := lat1 * PI() / 180.0;

lat2 := lat2 * PI() / 180.0;

a := POWER(SIN(dlat / 2), 2) + COS(lat1) * COS(lat2) * POWER(SIN(dlong / 2), 2);

c := 2 * ASIN(SQRT(a));

distance := earth_radius * c;

return distance;

end;

$$ language plpgsql;