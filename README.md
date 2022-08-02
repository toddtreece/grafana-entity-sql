### Grafana Setup 

```
git fetch && git checkout toddtreece/entity-grpc-server
```

`grafana.ini` config:
```ini
[feature_toggles]
grpcServer = true
topnav = true
panelTitleSearch = true
storage = true
export = true
dashboardsFromStorage = true
```

### Running the server

You will need to create a [service account](http://localhost:3000/org/serviceaccounts) with an `ADMIN` role, and generate a new token. Copy the token, and set the `GRAFANA_SERVICE_ACCOUNT_TOKEN` env var before running the next command.

```
go run . -token $GRAFANA_SERVICE_ACCOUNT_TOKEN
```

You should then be able to connect to it using the `mysql` cli:

```
mysql -h localhost --protocol tcp grafana
```

Example query that will run against the `devenv/dev-dashboards` path:

```
Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> select JSON_EXTRACT(payload, "$.title") as title from dashboards;
+------------------------------------------+
| title                                    |
+------------------------------------------+
| "Datasource tests - Postgres"            |
| "Datasource tests - Postgres (unittest)" |
+------------------------------------------+
2 rows in set (0.01 sec)
```