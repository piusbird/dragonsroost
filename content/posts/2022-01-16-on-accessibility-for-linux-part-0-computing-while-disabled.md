---
title: "On Accessibility for Linux Part 0 Computing while Disabled"
date: "2022-01-16"
---

Note a version of this post first appeared in my Mastodon Feed yesterday. This is an extended version with more detail. And is post one of my attempt at the [#100DaysToOffload](https://100daystooffload.com/) challenge

![](/assets/images/vsss.png)

## The Issue

I use a custom built screen reader-ish program to well read stuff on the computer.  Although screen reader is a bit of a misnomer. It' really more of a mutant hybrid between a screen reader, and a program meant to aide dislexics. I call it [vsss you can download it from my github](https://github.com/marnold/vsss) 

I will make a post on how to get it working in a non me sort of context later possibly tomrrow. In the meantime I'll say you need to rewrite vsss.conf.in for your system. So far as I'm aware Cepstral Swift is the only speech engine that supports the hooks needed for the fancy on screen graphics

When i upgraded to [#Fedora](https://tilde.zone/tags/Fedora) 35 on Thursday it stopped working. Not a problem in my code. I checked. Here's what went down, and why i'm so mad. Basically you have three options for text to speech on linux. First Espeak, Second Svox, (android tts) Third proprietary software synthesizer. And yes i know about the CMU stuff and hardware options. But for various reasons those aren't viable in my case

For various reasons i've used option 3 for the better part of 15 years now. And changing my computer voice now would be a huge adjustment. So the dirty little secret of most non-free speech synthesizers is they treat Linux/Unix as a third class platform. i.e most of the decently priced ones are still using OSS apis in 2021This hasn't been a problem as pulseaudio has this nice LD\_PRELOAD shim, that turns OSS apps into regular pulse clients.

You wouldn't think this would be a big deal for pipewire either; it is backward compatible with pulseaudio clients after all. Turns out it's not. But it turns out that for some reason that was not documented anywhere i could find [#Fedora](https://tilde.zone/tags/Fedora) dropped the shim for it in a recent update.

## Breaking Changes Strike Again

All the Changelog really says is that OSS, among other things is no longer supported. Without explaining why.  I suspect it's because almost no one uses OSS APIs anymore. ALSA has been around for 19 years now, and pulseaudio for ten. But I'm stuck with a binary blob compiled in 2012. Which from the tiny bit reverse engineering. I had to do for this project last saw major code changes, in 2007.

## A Convoluted and painful Journey

There is no technical reason why legacy OSS apps can't use the padsp shim to connect to a pipewire server. In fact I have this working. But in order to get it working I had to.

- Figure out that padsp had been removed from Fedora's pulseaudio package
- Attempt to revert the change in the source rpm.-- Watch that fail spectacularly
- Fish through upstream git to determine for a bit, reading the source code of the missing component determine that yes my theory was technically sound
- Uninstall Fedora's pulseaudio, try replacing that with upstream git build.
- Watch as my entire sound system explodes.
- Revert that
- Recompile pulseaudio again this time installing it under /opt -- play games with the linker so only the programs that actually need the replacement pulseaudio libs can see them
- IT WORKS

All this with diminished functionality in the reading software i depend on. All told this took 2 days 4 hrs and 21 minutes to figure out.

## The Takeaway

Moral of the story. An seemingly inconsequential change to you. May have catastrophic effect, on users with disabilities.

I was literally fired from my first job out of college four months in because of a Computer Accessibility issue I was unable to solve in a timely fashion. Recently I had to drop out of Graduate School for similar reasons.

FOSS has always had the potential to be the great equalizer bringing the digital revolution to the most marginalized and all that uplifting blah blah blah from my youth. And in my case it worked out I was able to totally replace a piece of software which ranges in price from $700 to $10,000 depending on vendor, feature set and so forth. With what is lets face it a radioactive shell script horror. From Ken Thomson's nightmares.

I was able to free myself from the lifetime of constant hardware and software upgrade costs that are often imposed on neurodivergent folks. In order to do basic things like read and write. Which is good, I doubt I would've completed college successfully without it.

But to a disabled person without my skillset the whole movement is a dead letter. Heck even if I was able to replicate my setup for someone. It is fragile as we have seen. We need to get better as a community at not breaking user space as Torvalds might say.

## Preemptive Troll Management

Now one might reasonably say "Why are you using Fedora if you need a stable platform", To which I say why should my disability exclude me from the latest and greatest, Pipewire can do amazing stuff and I was genuinely excited about it. If not for this seemingly random and unnecessary dropping of an essential tool for me, i would be quite happily playing with wireplumber and things if that nature.
