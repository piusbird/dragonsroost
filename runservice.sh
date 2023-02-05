#!/bin/bash

export SOCKPATH="/home/blog/http.sock"
export DSN="newblog:2mrE6CWih46D5mU@tcp(127.0.0.1:3306)/newblog?charset=utf8mb4&parseTime=True&loc=Local"
daemonize -o /home/blog/blog-startup.log -e /home/blog/blog-startup.log -p /home/blog/blog.pid -l /home/blog/blog.lock \
	-c /home/blog/dragonsroost -E DSN=$DSN -E SOCKPATH=$SOCKPATH /home/blog/dragonsroost/dragonsroost

chgrp www-data $SOCKPATH

chmod g+w $SOCKPATH
