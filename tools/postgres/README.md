# Postgres

This is a shared PostgreSQL instance for all apps

# Adding a new database

```
kubectl exec -it postgres-xxx --context homelab -- psql -h localhost -U xxx --password -p 5432 xxx

You should now be inside psql:

create database newdatabase;
grant all privileges on database newdatabase to xxx;
```
