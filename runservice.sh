#!/bin/bash

export SOCKPATH="/home/blog/http.sock"
daemonize -o /home/blog/blog-startup.log -e /home/blog/blog-startup.log -p /home/blog/blog.pid -l /home/blog/blog.lock \
	-c /home/blog/dragonsroost -E DSN=$DSN -E SOCKPATH=$SOCKPATH /home/blog/dragonsroost/dragonsroost

chgrp www-data $SOCKPATH

chmod g+w $SOCKPATH
