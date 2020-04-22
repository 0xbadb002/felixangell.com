#!/bin/sh
cd blogc; go build; mv blogc ../blogs/blogc; cd ../
cd blogs; ./blogc config.json