# Simple G-counter CRDT

A simple proof of concept about grow only counter crdt implemented in go lang. Purpose of this project is to learn about CRDTs and distributed programming concepts.

### How to run

1. cd into `server`
2. run `REGISTRY_URL="http://localhost:5000" PORT=3000 go run *.go`
3. You can start using the server independently
4. To add more replicas you can start another instance of the server
5. To connect the replicas and start merging the values between servers start the registry by cd'ing to `reqistry`
6. And start it `PORT=5000 go run *.go` command
