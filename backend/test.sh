set -e
export TESTMODE=true

for d in $(go list ./... | grep -v vendor); do
    go test -v $d
done
