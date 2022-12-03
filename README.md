# stress

### Postgres Replication
helm install postgres-replica bitnami/postgresql-ha --set postgresql.password=postgres
export POSTGRES_PASSWORD=$(kubectl get secret --namespace default postgres-replica-postgresql-ha-postgresql -o jsonpath="{.data.postgresql-password}" | base64 -d)
echo $POSTGRES_PASSWORD


kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

### K9S Metrics
kubectl patch deployment metrics-server -n kube-system --type 'json' -p '[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'

### ELK
docker build .
helm install logstash elastic/logstash -f values.yml

helm install kibana elastic/kibana 

helm install elasticsearch elastic/elasticsearch -f ./values.yaml

docker build . -t logstash-local
# saga-pattern-golang
# saga-pattern-golang
# saga-pattern-golang
