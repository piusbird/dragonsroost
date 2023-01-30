---
title: "Let's Create a 90s Web Programing Abomination Using Modern Technology"
date: "2022-05-03"
---

I make no secret of the fact I hate modern Web Development stacks. It doesn't matter what language or framework it is. I hate it all equally, I hate PHP even more. Flask sucks less then all the others. But it still sucks. Granted I still use modern web tools, money, and food dull the pain slightly. As does the opportunity to abuse sqlite in ways it's designers would find horrifying.

Why do I hate modern web stuff so much you might ask. Easy too much boilerplate, and deployment is a nightmare. The programs I want to expose to the web are usually silly little one offs. Most of what i want to do web doesn't even merit a database connection. Let alone a full model view controller stack, containers, an ingress controller uWSGI and all the other goop, that one simply must have.

## Once Upon a Time

It used to be simple to write web applications here's an example

#include <stdlib.h>
#include <stdio.h>
     
    int
    main(void)
    {
     
    	puts("Status: 200 OK\\r");
    	puts("Content-Type: text/html\\r");
    	puts("\\r");
    	puts("Goodbye Cruel World!");
    	return(EXIT\_SUCCESS);
    }

Ok What's going on here. Well It's easy m'kay. A web program was just a normal program which prints out http headers followed by a blank line, followed by whatever generated content you wanted. HTTP info was stored in well documented environment variables, and the program had to return a success code on exit. To use just compile, upload and [presto](https://piusbird.space/~matt/goodbye.cgi). No docker or proxy passing required. This was called the Common Gateway Interface, and it was the backbone of dynamic content on the web for 15+ years. Heck Most modern web stacks just build or elaborate on this simple protocol

## Why did things change?

Scalability, Security, and so forth. CGI as originally implemented spawned a separate process for each request sent to the server. Which could bog down busy sites quite easily. On the security side. Well I won't go into it here but [this article](https://antonyt.com/blog/2020-03-27/exploiting-cgi-scripts-with-shellshock) is quite nice if you want to look into it yourself.

But as I said above I hate the modern web, Most of it is overkill for the stuff i do. I know the security risks involved in using classic methods, and scalability concerns for the sites hosted here is a problem I'd love to have. So the question then becomes

## How to get Classic CGI working on a modern webserver, by which I mean Nginx

![](/assets/images/218483.png)

<Troll> Why not Just switch to Apache CGI works great over there

I'm not switching back to Apache because I don't like the memory hogging tendencies or the configuration file format is just bad m'kay.

<Troll> What is Nginx, I thought everyone used Apache

Nginx is a webserver sorta like Apache. In fact it is currently the most deployed webserver on the internet, it surpassed Apache in that role in about 2016 as I recall. Here's the latest survey data I could find.

![](/assets/images/web-server-usage.png)

Webserver data

Nginx is much faster, and much less of a resource hog than Apache. But for present purposes there's a problem Nginx has no ability to serve dynamic content on it's own. Meaning no cgi support, no php support no nothing. What Nginx can do is pass http requests to so called application servers sitting behind it. Either through a protocol called FastCGI/WSGI, or a plain old reverse proxy. When it gets a result, it does some quick header rewriting and displays that to the user. This saves resources, has security benefits and also allows you to scale and load balance application, should you become the next Facespace or something. All this is great, and most people love it like 80% of the time the other 20% being used to curse out the inevitable 502 Bad Gateway Errors which you will get if you try to do some of the more advanced Nginx tricks.

From this description it should seem obvious what we have to do. Find an application server for use with Nginx that runs old style CGIs. Configure it, and profit.

Sort of like this

location ~ ^/(~|u/)(?<user>\[\\w-\]+)(?<user\_uri>/.\*)?$ {
          alias /home/$user/public\_html$user\_uri;
          disable\_symlinks if\_not\_owner;
          autoindex on;
          
          location ~ (\\.cgi|\\.py|\\.sh|\\.pl|\\.lua|\\/cgi-bin)$ {
             gzip off; 
             include fastcgi\_params; 
             fastcgi\_pass unix:/run/slowcgi.sock; 
             fastcgi\_param SCRIPT\_FILENAME $request\_filename;
          }
          
        }

## Enter SlowCGI

From my configuration snippet above You'd think this would be simple, but it turns out until about whenever this post goes up it wasn't. There's not much reliable documentation on how to do it, and what exists is either old, or very distro specific. So I decided to work this out on my own with a little help from IRC as usual thanks xwindows.

So our first hurdle is as mentioned the application server piece. Turns out there are two appservers that allow you to run legacy CGI applications. One is known as fcgiwrap, and one is called SlowCGI

Both have problems as it turns out fcgiwrap is nearly unsupported, and after two hours of fiddling I couldn't get it to work on Ubuntu 20.04. Although it works great under Fedora so there's that. :P.

SlowCGI is actively supported by the folks at OpenBSD, but is alas an openBSD exclusive application. So to make it work i needed yo port it to Linux. Which turned out to be trivial, heck most of the work was already done WAY back in 2018.

I just needed to sync the code with upstream and make a few changes to Makefiles systemd units and so forth.

It's over on [sourcehut](https://git.sr.ht/~marnold128/slowcgi-port)

The most painful part of the `port` was figuring out the Ubuntu/Debian used LDLIBS instead of LDFLAGS. Which took about an hour of googling to figure out.

The last bit was making the systemd unit file work on Ubuntu, and configuring Nginx to use it. Which you can see above

## Closing Thoughts

Be careful with this legacy CGIs have security implications beyond just Shellshock. In the default configuration anything which the webserver has permission to read/write can also be read/written by the CGI program.

Also worth noting is the fact that SlowCGI is less tolerant of badly coded scripts. Be sure to send at least a Content-Type header, and the all important blank line at the end of headers, or you'll get the dreaded 502 error, with only cryptic messages in the log to guide you.

If i ever follow up this post I will include information about how to use BubbleWrap to mitigate some of the security concerns.

Meanwhile Embrace the Joy of Linux everybody

/Matt
