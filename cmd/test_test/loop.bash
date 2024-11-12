#!/home/das/.nix-profile/bin/bash

start=1
end=100

for ((i = ${start}; i <= ${end}; i++)); do
    echo "--- Iteration #${i}: $(date) ---"
    time go test -v
done