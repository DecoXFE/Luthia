-- Drop tables in reverse order of dependencies (children first)
DROP TABLE IF EXISTS job_events;
DROP TABLE IF EXISTS workers;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS workflows;
