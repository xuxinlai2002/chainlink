#!/bin/bash

echo "start run postgres db"
docker run --name cl-postgres -e POSTGRES_PASSWORD=pdb0815pdb0815pdb0815pdb0815 -p 5432:5432 -d postgres:15.6
