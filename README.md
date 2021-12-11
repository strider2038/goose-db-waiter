# Database waiter (based on goose migrations)

Database waiter can be used as an init container for Kubernetes deployment that must wait until database migration is completed.

Example of usage

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  # metadata specification
spec:
  # deployment specification
  template:
    # metadata specification
    spec:
      containers:
        # microservice specification
      initContainers:
        - name: db-waiter
          image: strider2038/goose-db-waiter
          imagePullPolicy: IfNotPresent
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: config_name
                  key: DATABASE_URL
            - name: MIGRATION_VERSION
              valueFrom:
                configMapKeyRef:
                  name: config_name
                  key: MIGRATION_VERSION
            - name: MIGRATIONS_TABLE # optional, default is "goose_db_version"
              valueFrom:
                configMapKeyRef:
                  name: config_name
                  key: MIGRATIONS_TABLE
```
