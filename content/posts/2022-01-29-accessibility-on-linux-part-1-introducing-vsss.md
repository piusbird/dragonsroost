---
title: "Accessibility On Linux Part 1: Introducing vsss"
date: "2022-01-29"
---

Hello everyone and welcome back to another edition of Piusbird attempts to build his portfolio.  Also known as somebody please hire me please; I'm competent I swear.  In a previous post I said I would outline how my custom screen reader worked, and more importantly how to get it working in a non-me context. Well ladies and gentlemen the time has come.  Here's a quick primer on the Very Simple Speech Service.

## Design Explained

  Upon cloning the code from git you might be tempted to question my sanity.  It is after all a polyglot program composed of one D-Bus-based micro-demon, a large and seemingly complicated shell script and an optional screen handling routine written in C..  But I promise you this Goldberg-esque madness was all rationally designed.

## The Method To My Madness

I needed a speech system that could be deployed on any Linux or UNIX system with X11. Which depended only on those things which I could commonly find on the systems I was using at the time.  With minimal if any required additional package installations and especially no additional python module installations at all. Thus it is written mostly in bash. The second design constraint that dictated I write most of this in shell script was it needed to be adaptable to any environment that I came across as I had no guarantee of root access to any system.  Thus it is possible by simply changing a couple of variables to do without the Python micro service or the fancy screen handling stuff.

The last design challenge which made my odd choice of language reasonable was I only had 12 hours to get the first version out the door. As I remember it had something to do with finals week of 2013  

This means that the version on GitHub is specifically configured for my setup and what we will be doing in the remainder of this post is adapting it to yours or attempting to at least.

Before we get started in earnest I should mention that when I say X11 I really do mean X11.  Logically there is no reason why it shouldn't work on Wayland. I get odd warnings from weird places when I have tried it and since I see no reason to use Wayland yet I have not looked into it further.

## To The Code

I mentioned earlier that my code was highly adaptable for any environment. And while that is true there is one hard and fast dependency. This is of course a software speech synthesizer. It can be any one you'd like, however festival or espeak-ng are commonly installed by most distributions. The remainder of this guide assumes we will be using espeak. Simply because that's what I have in the VM I'm testing this with.

The first file you'll need to modify is called vsss.conf.in which looks like.

## vsss.conf.in

```
VOX="Callie" 
audio_bckend="padsp" 
rate=200
spkedit="pluma" 
QT_SELECT=4; 
export QT_SELECT

PIPE_COLOR="1;33;44m"
export PIPE_COLOR
speak_bckend() {

    if [ -f /tmp/vsss.lock ]
    then
	echo "Speech output is currently in use"
	return 0
    else
	touch /tmp/vsss.lock
    fi
    $audio_bckend swift -n $VOX -p "speech/rate=$rate" -f $1 -m text -t | colorize-pipe 
    rm /tmp/vsss.lock
    return 1

}
```

This File has one and only one job. To define the speak\_bckend function and any supporting variables, or other functions it may need. This function is what actually does the speech synthesis, and takes one parameter. A file name which contains the text to be spoken. In my setup this function depends heavily on Cepstral Swift, and it's quirks. Let's change it to make it use espeak.

```
# Espeak Version 
rate=200
spkedit="nano" 
speak_bckend() {

    if [ -f /tmp/vsss.lock ]
    then
	echo "Speech output is currently in use"
	return 0
    else
	touch /tmp/vsss.lock
    fi
    espeak -f $1 -s $rate
    rm /tmp/vsss.lock
    return 1

}
```

Note a couple of things here. First it is best practice to implement a lockfile mechanism, before allowing the speech synthesizer to execute. Unless your a fan of symphonic chaos of course. Also note i set the speech rate to 200 words per minute. To those unpracticed with text to speech this can seem almost incomprehensibly fast. But keep in mind this is actually 100 words a minute slower than your average adult reading with their eyeballs. If your having trouble understanding the computer slow it down to about 165.

I've committed the espeak config file to github so you should be able to just copy it over top of mine, and done. And please note I'm always open to merging pull requests for more back ends.

## Rarer Modifications

The next two modifications I will show only apply if you don't have QT or dbus installed. In this case you will need to comment out line 9 of vsss,

And finally change line 20 in vsss\_cmd.sh to fetch your primary clipboard without using my dbus service something like

```
xsel -b  
```

Should work fine

## Ready, Set Go!

Assuming everything went according to plan and you are running the latest git, from Friday 28th January 2022. You should now be able to run ./vsss and be greeted with something like

Very Simple Speech Shell  
Version 0.3.4+test

\>>>

## Conclusion

I hope this was enough to get the prospective user of my very odd reading software through the process of setting it up/ The last post in the series will cover actually using it to get work done.
