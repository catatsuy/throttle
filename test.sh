#!/bin/zsh

for i in {1..10000}; do
  sleep $(((RANDOM%3+1.0)/10))
  echo "Welcome ${i} times"
done
