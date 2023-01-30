---
title: "New Year New Proxy"
date: "2023-01-06"
---

![](/assets/images/out-0.png)

AI Image Prompt: the legend of zelda breath of the wild, 1girl, purple hair, laptop computers, happiness

[BlueProxy Git](https://git.piusbird.space/dockerfiles.git/)

This was not a project I hoped to be announcing today. I was hoping to announce the start of my solo RPG adventuring platform. I spent my November and December fighting the mental health system. So I had to scale back my more ambitious projects. Nevertheless I managed to create something that I’ve always wanted to create. To close out 2022 and began 2023 in the right way.  

This project is what I call an accessibility proxy. An accessibility proxy is a website which takes the contents of another website of the user’s choice and rebuilds the visual elements to be more usable by people with print related disabilities. Print related disability is an umbrella term covering anything from blindness to dyslexia and much else in between. Basically anything which affects a person’s ability to read standard text can be classified as print related disability.

## Technical Overview

![](/assets/images/Screenshot-from-2023-01-06-15-10-13.png)

Netsurf showing the landing page for my proxy program. It Welcomes you and then asks for a url and user agent.

I grew frustrated with not being able to read things on the Internet. So I had a “hold my beer” moment over new years. I ended up coding the entire thing in 72 hours or less. It hasn’t hit the forges yet but you can find it on my personal Git server.

The heavy lifting in this piece of code is done by a fork of a piece of software called “miniwebproxy” by openBSD’s Ted Unangst.

This works fine for now, but I eventually hope to replace it with an engine based on Mozilla’s Readability, which can handle even more sophisticated modern websites. I added selectible User Agent support, and Header Sanitization.

The next release will see support for converting to plaintext. The final feature of note in the first release is the web frontend is completely JavaScript free, and it was designed to work with older or underpowered browsers e.g. Dillo, and NetSurf. It might even work with Netscape 2.0 but I haven’t tested this as yet. Note that I don’t plan on running a public instance of this as this piece of software has shall we say a tendency to bypass paywalls, and make “Website DRM” magically disappear. My instance is hiding behind certificate-based authentication. So as not to get blocked. However I have open sourced the project as usual so that if you want to start your own instance you can.

![](/assets/images/Screenshot-from-2023-01-06-15-14-56.png)

The final proxy render

## Bugs of note:

CloudFlare tends to be a pain in the ass when using exotic user agents. Thus a standard Firefox user agent is provided. I will add more normalization features soon. The Dockerfile doesn’t seem to want to build on Oracle cloud Free Tier and I have no idea why. It correctly builds everywhere else I’ve tried.

Happy reading and enjoy if you have any feature suggestions pull requests etc. my inbox is always open.
