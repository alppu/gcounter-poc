# Simple G-counter CRDT

A simple proof of concept about grow only counter crdt implemented in go lang. Purpose of this project is to learn about CRDTs and distributed programming concepts.

### How to run

1. cd into `server`
2. run `REGISTRY_URL="http://localhost:5000" PORT=3000 go run *.go`
3. You can start using the server independently
4. To add more replicas you can start another instance of the server
5. To connect the replicas and start merging the values between servers start the registry by cd'ing to `reqistry`
6. And start it `PORT=5000 go run *.go` command

### Included features

- Each replica can be used independently of each other and without registry
- Each replica tracks everyone's states so that there is no data loss even if some of the replicas go down
- Automatic merging of other replicas with go-routiens when incrementing values
- Basic service registry to store and provide other replicas addresses
- Server caches other services addresses so even if the service registry goes down replicas can sync data between each other
- Server implements circuit breaker to remove dead services from cache
