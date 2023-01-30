---
title: "What have we lost?"
date: "2022-09-13"
---

## Video for this post

https://www.youtube.com/watch?v=AIhy0T7Q48Y

A modern day idol of mine Xe Iaso recently, wrote a post on the end of [Heroku Free Tier](https://xeiaso.net/blog/rip-heroku). In which she stated that, this service was vital for her getting the career that she has. It struck me that in order for what she says to be true, which I assume it is she must be at least five to ten years younger then me. She says she's watched as doorway after doorway into tech has been firmly shut and locked. That is truer than I think she knows. I wanted to expand upon this idea a little.

## In My Day

Back when I cut my teeth in tech, free hosting with PHP, and MySql was plentiful and ad supported. If you wanted more you could always install Linux on an old Pentium II, do some port forwarding and dynamic dns magic and by god you had a website. Or if you were not inclined to Linux you could run any of the myriad of WAMP stacks in a box such as EasyPHP. If you were truly desperate, there were services like geocities. Which gave you a little less but it got your foot in the door.

The entire concept of `free compute power` was not something that existed back in 2004. Learning Unix or DNS or PHP was hard. But and this is the key difference. Everyone starting out had to learn the same skills. When we all got our first jobs those skills were what the interviewers looked for. Point is every middle class white kid started from the same level. This state of affairs endured for about a decade. However, things began to change in 2012.

It started with Linux distros slowly dropping support for older hardware. To the point where today you need x86\_64v2, introduced in 2013 to even boot in some popular distros today. The problem with this is, that the useful life of a machine is now way beyond ten years. Such thinking only considers the paying customer, and not the poor kid just trying to find a start

## Now

Some residential ISPs now no longer allow port forwarding either. And even if these two states of affairs were corrected. Modern web technology is now so tied up with "The Cloud", that it would be of limited utility anyway. Is there even a path for anyone to learn how a website is really built these days? I have no idea.

I never thought i would be saying this but Thank God for Oracle Free Tier, and fly.io. For however long they endure.

## Only 25% of the problem

But the problem of resources is only half, or less of the issue I see today. We are now so abstracted, and hyper specialized that it is hard to know just what constitutes, "General Tech skills". I think a few examples would be illustrative here.

Docker now dominates the world, back in 2012 when it started out I was fundamentally a distribution package maintainer. Even in my professional life and freelance work I did things like build packages for legacy software on Ubuntu 12.04. 2012 era docker was ill suited to that workload so i ignored it. Even though i understand what Docker does and how to use it. I think it's overkill for the projects I have now. Because of this recruiters looking for Docker or k8s `experience`, pass me by. So nearly every recruiter these days passes me by. This is only the tip of the iceberg.

Do i specialize in Go, Python or Rust these days? They seem to be mutually exclusive in the minds of many an HR department or Recruiting firm. If I had to define a course of study to intentionally take a green as grass teenager, from newbie to pro along the lines of what i had in my childhood. I doubt i could figure one out, without also irrevocably shaping that teenager's future career prospects. But the iceberg goes still deeper.

## Story Time

Have you ever heard of linktree? I won't link it here because I am fundamentally opposed to their entire existence. What they give you is essentially a place to put links to all your social media apps. So you can give that URL to all your friends and followers. So they can follow you everywhere. If your saying "There's a platform for lists of hyperlinks". The answer More than one.

The reason i know about this is one of my younger friends suggested i get one, because the number of things I am doing these days is only growing. While the concept itself was sound. I wanted to host it on my own domain. I looked for a self hosted version. But the popular one involved github-actions visual studio code and a bunch of other overkill that made my brain melt into my shoes. I wanted to write it in Markdown. Compile it to html, scp that to my site and have done with it. So I wrote a Makefile. It wasn't till i looked at the output when it wasn't styling correctly. That i saw it was just an unordered list of anchor tags. No different from my first hand coded website of 20 years ago.

So what did I do? i fired up nano and made an unordered list of anchor tags warped that in a div for easier styling. Did some quick css et viola

When i showed this to the younger friend who had instigated these shenanigans. Far from being impressed, she said "Your doing MySpace S\*t, why not go with something easier". She was genuinely confused when i explained that html was still how all websites work.

## The Lost Tools of Learning

The moral of the story is that the early internet services sort of forced you to learn how the internet itself worked in order to fully express yourself.

You had to know about things like html, file paths, midi, and a whole host of other things to even begin to have a cool geocites site. That is no longer true these days. Now platforms take pains to hide the inner workings from the user as much as possible. Some even going so far as to lock people in to their own templating system. When the bother to support user styles at all.

I remember remarking on this aspect to a friend of mine when Facebook opened up. She accused me of being a `tech elitist`. Maybe i was back then. But now i think that having at least a basic understanding of these fundamentals, is a key on ramp into Tech itself. With platforms taking pains to hide the inner works of the web from everybody. It's a wonder that anybody finds their way in the door these days.

So what do you think? Has capitalism destroyed the potential of this technology as we once knew it? Let me know in the usual places

Pius
