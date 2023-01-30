---
title: "More 90s Web Programing Abominations"
date: "2022-05-08"
---

In our last adventure we got SlowCGI working under Ubuntu. That was great. But one must wonder Why? Well I did it because I wanted to learn how to make web applications in pure C. Both to swat up on my rusty skills in C, and for no other reason than I want to learn the skill. You see about 10 years ago I interned at a mid sized enterprise cloud provider with roots going back to the late 90s. During which they produced some of the first web applications for the Federal and State Governments. They were all done in pure ANSI C89. Why Because none of the web technologies of those days think Perl, and PHP 2.0, ColdFusion. Could work at the scale required, this was common in the early days, from what I gather. Yahoo Stores for example was written in LISP for example.

## A bucket list item

My mind was blown when one of my supervisors showed me the contents of their old CVS repos. And it's been on my `programing bucket list` to learn how to do webapps `the old fashioned way`. Thus my port of SlowCGI to make it possible. And recently i've completed my first dynamic web toy in C. It's [here](https://piusbird.space/~matt/ascii.cgi) source is [here](https://piusbird.space/~matt/random/squarecgi.html). Output example [here](https://piusbird.space/~matt/ascii.cgi?v=Piusbird)

## Notes on CGI in C

These are my notes on the pitfalls of it. First off you really want CFLAGS set to `-Wall -Werror` when doing this sort of thing because the compiler will allow you to do some pretty stupid stuff that will result in segfaults otherwise. But that's true of almost every C program i find. For those not in the know `-Wall -Werror` is a mode in modern C compilers which will warn you about code that is potentially problematic. And then treat those warnings as compiler errors. Think of it like training wheels for the C compiler. Granted not all warnings are valid, especially in older code. But for those returning to C after 10 years it's a god send

I didn't have this turned on my first try. So I spent an good two hours trying to figure out why segfaults? They were happening despite my use of only safe string functions and very early on in the program. With -Wall -Werror enabled I found out my strings weren't being initialized properly, because of C's order of operations.

This needs a better method of parsing form data. I use a rather bone headed method, that while `safe` from a memory access perspective. It can't decode spaces in the form data for example.

We also really need some sort of template engine, Editing Strings in a C file, every time the html needs to change sucks

## Conclusion

I hope someone gets something out of my silly little hobby project.

Till next time Embrace the Joy of Linux everybody
