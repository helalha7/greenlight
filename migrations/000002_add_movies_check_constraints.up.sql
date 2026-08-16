ALTER TABLE movies
ADD CONSTRAINT movies_runtime_check 
CHECK (
    runtime >= 0
);

ALTER TABLE movies
ADD CONSTRAINT movies_year_check 
CHECK (
    year >= 1888
);