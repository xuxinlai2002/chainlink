The files belonging to this database system will be owned by user "postgres".
This user must also own the server process.

The database cluster will be initialized with locale "en_US.utf8".
The default database encoding has accordingly been set to "UTF8".
The default text search configuration will be set to "english".

Data page checksums are disabled.

fixing permissions on existing directory /var/lib/postgresql/data ... ok
creating subdirectories ... ok
selecting dynamic shared memory implementation ... posix
selecting default max_connections ... 100
selecting default shared_buffers ... 128MB
selecting default time zone ... Etc/UTC
creating configuration files ... ok
running bootstrap script ... ok
performing post-bootstrap initialization ... ok
syncing data to disk ... ok


Success. You can now start the database server using:

    pg_ctl -D /var/lib/postgresql/data -l logfile start

waiting for server to start....2024-08-17 04:16:05.989 UTC [72] LOG:  starting PostgreSQL 15.6 (Debian 15.6-1.pgdg120+2) on aarch64-unknown-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
2024-08-17 04:16:05.990 UTC [72] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
2024-08-17 04:16:05.998 UTC [81] LOG:  database system was shut down at 2024-08-17 04:16:05 UTC
2024-08-17 04:16:06.006 UTC [72] LOG:  database system is ready to accept connections
 done
server started
CREATE DATABASE


/usr/local/bin/docker-entrypoint.sh: ignoring /docker-entrypoint-initdb.d/*

waiting for server to shut down...2024-08-17 04:16:06.217 UTC [72] LOG:  received fast shutdown request
.2024-08-17 04:16:06.218 UTC [72] LOG:  aborting any active transactions
2024-08-17 04:16:06.223 UTC [72] LOG:  background worker "logical replication launcher" (PID 84) exited with exit code 1
2024-08-17 04:16:06.225 UTC [75] LOG:  shutting down
2024-08-17 04:16:06.226 UTC [75] LOG:  checkpoint starting: shutdown immediate
2024-08-17 04:16:06.273 UTC [75] LOG:  checkpoint complete: wrote 918 buffers (5.6%); 0 WAL file(s) added, 0 removed, 0 recycled; write=0.022 s, sync=0.021 s, total=0.048 s; sync files=301, longest=0.011 s, average=0.001 s; distance=4223 kB, estimate=4223 kB
2024-08-17 04:16:06.284 UTC [72] LOG:  database system is shut down
 done
server stopped

PostgreSQL init process complete; ready for start up.

