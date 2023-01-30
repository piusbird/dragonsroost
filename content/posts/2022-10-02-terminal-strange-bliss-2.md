---
title: "Terminal Strange Love"
date: "2022-10-02"
---

  
Or, How I learned to stop worrying, and learned to love the kludge

![](/assets/images/Screenshot-from-2022-10-01-19-26-43-1024x549.png)

## Introduction: An Unlikely Wordsmith  

> The objective reality of the situation seems to disprove my self assessment. And in spite of my lack of self confidence in this field, I do enjoy slinging words about.
> 
> Me

Everyone tells me, that I am a good writer. Well nearly everyone. My sixth grade English teacher was an absolute witch bent on destroying the self esteem of impressionable 12 year olds, and doesn't count. And there was that low grade on the GED, but we can put that down to a bad day. Point is I have a 90 average on all English/Creative Writing assignments in my college years. More importantly I have been published in print and online dozens of times over the years, beginning when I was 9 years old.

Although I consider myself a two bit hack with no voice, and a style which would make your average third grade teacher scream in agony. The objective reality of the situation seems to disprove my self assessment. And in spite of my lack of self confidence in this field, I do enjoy slinging words about.

## Striving  

However there is a problem. My neurodivergency and disability make writing one of the hardest tasks I regularly undertake. It is both physically and mentally exhausting to attempt to write even a simple three page academic essay. Mentally because I have ADHD, and thus a lack of working memory, and high distractibility. And the Cerebral Palsy makes both eye muscle strain and hand, and finger pain an issue.

To be clear this applies to all computer interactions to some degree or other. These problems take on there most severe forms when I attempt to write. Also games like Minecraft are nearly impossible for me to enjoy. But that is a different blog post entirely.

This has been a persistent and annoying problem throughout my life. And various Teachers/Occupational Therapists/Technologists etc, have tried dozens of solutions through the years. From forcing me to do as much writing as possible in hopes that the problem would go away eventually.

To various software solutions such as Don Johnson's CoWriter to Dragon Naturally Speaking, and many others. As you might expect the exposure method was almost torture, and the software solutions have had mixed effect.

Here's a hint though in my last two years of university I had to get a grant to pay an actual human being to type my papers which I either recorded to MP3, or dictated live. This problem has limited my writing ambitions quite severely. Needless to say there are no grants for Game of Thrones/Harry Potter FanFiction authors, and writing that kind of thing with Dragon is nearly impossible. So I have continued the search for a usable writing tool set and work flow, which works for me. After many years and hundreds of attempts I think I may have found my best solution yet.

In the remainder of this post I will outline the Engineering steps I took to arrive at this set up. To aide other disabled people in their pursuit of appropriate assistive technology.

## Learning From Failure

The first step in any engineering task, weather it be building a technological system or a bridge, is to identify why what you currently have isn't working. In some cases this may seem obvious. You need a new bridge because the old one fell down, for example. But the important question is why did your old bridge fall down? Or in our case Why have nearly all of the attempts to solve my writing issue eventually failed? My answer to this question is two fold. First the attempts by so called \`Professional Experts\` to solve this problem, were too narrow in scope. For example the Teachers who assigned the exposure exercise were seeing my distractability as a function of immaturity, and weren't necessarily taking neurological, or physical factors into account.

Likewise the PT/OT people were only thinking of muscle weakness. They weren't looking at how black text on a white screen, might not be as high contrast as needed to make eye strain less likely over time. I'm not sure if a vision expert ever commented on this issue. But the point is each expert was limited to seeing only the aspect of the problem which happened to coincide with their own field of expertise.

In a lot of ways the worst offenders here were the \`Technology People\` were even worse. With two exceptions they never worked with me directly, and only applied pre defined solutions to surface level problems.

To be clear I’m not trying to denigrate any of the people who tried their best throughout the years. I am simply pointing out that none of them had as full a view of the problem as my 24 years of experience writing has afforded me.

## Defining the Problem

This is an important lesson for all disabled persons. You are the expert on You. Don't be afraid to assert this. So what is this writing problem I speak of exactly. Eventually I broke it down into three parts.

A. Hand and eye pain  
B. Distractability.  
C. Speed of production.

These three aspects are mutually self reinforcing. For example it's hard to stay focused on the piece at hand when your hands are throbbing. And this slows things down which leads to more opportunity for distraction, said distraction slows things down yet again and pretty soon your only three words in to 5000 word term paper, and you're five hours from deadline.

## Typical solutions

The typical solution has focused on reducing or eliminating, my pain through solutions like word prediction and voice recognition.

As evidenced by this post that approach has had a profoundly mixed record of success. However in my latest stab at the problem I decided to attack the Distractability angle first.

For the simple reason that I have a lot of stuff to write in the next six months and I reason that if I can keep the words I wish to write in my head and get at least a standard page,(about 250-300 words) done on one project per day. Then by April I will several thousand more words than I do currently and will have met most if not all of the deadlines both self imposed, and otherwise. And in the process of finding a solution to the ADHD issues.

I think i may have found a system that addresses all three legs of the triangle at once. In fact i would go so far as to call it almost perfect. At least in my testing so far. However in order to get this far I had to define precisely what a viable solution would look like. In software engineering terms we call the definition of success Acceptance Criteria.

> The key to success is basically defining what success looks like before you start.

## My Acceptance Criteria

For me a successful system has several properties. We begin by defining what I call the environmental characteristics of the system. What hardware or software will this run on? How will it interact with the cloud? So on and so forth First and foremost it must be able to run on freely available, open source UNIX operating systems. I have used Linux and FreeBSD almost exclusively for the past 20 years of my life. I won’t go into the economic and ethical reasons for this. As you my primary audience probably knows all of these chapter and verse. Second and repeatedly any solution must be able to run on low-end hardware.

My newest computer was built in the year 2014, one would think this wouldn’t be a problem for word processing even heavyweight stuff. However outsiders are often surprised at how resource intensive assistive technology systems can be.

Thirdly any successful system should be able to function with minimal to zero cloud computing infrastructure. This is a critical point almost as critical as the first two.

## A short rant about modern software

Modern voice recognition software that transmits all spoken audio from your computer to the cloud service which sends back the text. This has several advantages for the average consumer. It eliminates were severely curtails the often laborious process of training the voice recognizer on an individual voice, and consequently dramatically increases both speed and accuracy of recognition. For a serious writer who is disabled however this seeming advantage turns into a disadvantage.

Because the voice recognizer isn’t trained on an individual voice the customization of local voice-recognition is lost. I said before writing fan fiction of a Game of Thrones or Harry Potter is nearly impossible with commercial voice recognition software. In recent versions 12 and above this is true.

But older versions will work just fine provided you’re willing to put in the training time.

But I digress point is all system critical components such as voice recognition spell check, and grammar and stylistic suggestions must be able to run disconnected from the Internet.

Non-critical features like file synchronization may use the cloud. However I decided that if this occurred I should own the infrastructure or leave town infrastructure portability for these features.

## The Open Source Puzzle

Notice what I didn’t place in the environmental criteria section. I have made no mention of the system being composed of open-source software. While this would be my preference,

I also know enough about the open source accessibility ecosystem to realize that placing this requirement on myself would be self sabotaging. Projects that benefit people with less common disabilities rarely get the developer time they deserve. In addition of this there are often patents concern with open-source implementation of the common accessibility systems.

So I am fine with running things in Virtual Machines or under an emulation layer such as Wine. So long as the system requirements are low enough to make that efficient. In practice this means obtaining software from the Internet Archive, off of eBay or from other sources. Now that we have define our environmental criteria and component selection model we can move on to what accessibility properties the system must have.

## Accessibility Features

Under accessibility properties I first want a system that minimizes distraction. This is a broad and nebulous statement so let’s see if we can narrow it down. To me a distraction if anything that triggers my ADHD, to the extent that I forget what I’m writing. This can be a variety of different things from the squiggly lines that modern word processors place under words and phrases to inform you about spelling and grammar errors, animations and other such modern user interface glop. Generally anything that draws my eyes away from the words I am writing is a bad thing.

Also among the key features of a successful writing system is the ability to use a high contrast or yellow on blue color scheme to minimize eyestrain.

But perhaps what I consider to be the most critical feature of this system is the ability to display a document on screen in a magnified view without affecting the fonts and styles of the finished document.

If you’ve ever had to deal with a teacher who would only give a their essay lengths as a page count and then accuse you of cheating when you wrote it in Comic Sans 14. Or worse yet had to hand in a formal report to your boss with the aforesaid infamous font set you will know how critical a feature this is.

There are a few other design criteria which I used when pulling together this new system of mine. But I will not bore the reader with all of the major design decisions I made in the course of the project.

I hope that this section has given you a flavor of what designing a system is like.

The key to success is basically defining what success looks like before you start.

## Implementation Details

Now that we have defined what success looks like I can tell you all about the system I came up with. Which as I have the best iteration of the writing system to date. That is if you believe my productivity metrics for this week. 5000 words written in four days. This easily beats the previous iteration of my writing system which can only allow me to produce 1500-2000 words in that same time span.

A 250% increase in productivity seems pretty good to me.

Solutions like PanDoc and markdown, do not work for me on larger projects. Anything over 500 words and I spend more time fiddling with the markdown than actually writing. I have also found spell check and grammar correction to be a pain when writing and markdown or similar solutions.

There is one feature I like about writing and markdown however. You can write in a terminal-based editor such as vim or nano. Which pretty much gives you the freedom to edit the document using any monospace font and size that you wish to use. And no one reading the finished product will know that you were writing in Comic Sans Code with ligatures. Which as stated above, is a key success criteria. So we have a feature that we know we like and also a shortcoming that is possibly a showstopper.

Now we must ask the question are there similar technologies with this feature available that lack the undesired behavior.

## Retro Futurism

![](/assets/images/Screenshot-at-2022-09-30-18-55-39-1024x549.png)

This is where terminal word processors came to mind. For those unfamiliar with computer history, terminal word processors existed in an era from about 1979 to 1994 when the computers mostly lacked the graphical capabilities that we take for granted today. In those days most of the graphical capabilities were located in the printer. Terminal word processors provided an easy way of controlling a printer’s graphical capabilities for the average user.

The idea was on the screen would not be exactly what would appear on the page but an approximation of it; which was designed to let you predict with some accuracy what the printer would do. Without necessarily forcing to edit the document in the way it would appear on the page. I knew about this long abandoned branch of the tech tree of course.

There are science fiction and fantasy authors, including no less an authority then George R.R. Martin who swear by WordStar, and all the best lawyers I know will praise him WordPerfect to the skies and back. However until recently these steps needed to get a terminal word processor working on a modern system were so insane and wonky, that I just didn’t bother.

## Almost Perfect

And that was before one took into account the fact that I need voice-recognition sport on larger projects. Which implicitly means being able to interchange documents with a semi-modern platform.

Enter Tavis Ormandy who managed through a feat of digital archaeology to dig up the last WordPerfect for UNIX terminal release version 8.0.076 built for Red Hat version 6 released in 1998. He then proceeded to reverse engineer and patch the binaries so they would work on modern systems. As soon as I saw this come across my feed I knew I had to at least try rebuilding the writing system around this hidden gem. Two weeks later I had success.

## WordPerfect configuration

Through a lot of trial and error I discovered that best way to use WordPerfect was through an xterm, with a monospace font and 18 points.

With the foreground color set to yellow and the background color set to Navy as shown in the images. This minimizes eye fatigue as much as is possible. Meaning I am more accurate when it comes to missing words and things like that Rather than directly print to a postscript file as Tavis shows in his example, but rather I use printer-driver-cups-pdf. With WordPerfect's GhostScript driver to produce pdf files of documents.

> The TL:DR version of which goes something like. The potential of technology for the disabled will only be realized when, we stop kow towing to an industrial complex which doesn't know us, and is only interested in the money they can get out of us

## The Last Piece

This still leaves one problem however. I need voice recognition projects longer than 1000 or so words, after the 1000 word mark I start having trouble holding large chunks of the piece in my head, and remembering where I am going with any given point. When you add motor planning keystrokes.

My little ADHD brain quickly becomes overwhelmed. But how does that work you might ask. Sure it's easy to set up a low end XP VM, install Dragon 9 (the last good version in my opinion) and go to town. XP even has High Contrast built in. WordPerfect for UNIX however still saves it's documents in the 30 year old, WordPerfect 5.2 format. The code for the Unix/Linux version was forked off of 5.2 and never had a rebase before it was abandoned.

Ordinarily this wouldn't be a problem. Because LibreOffice is awesome. It can not only open but save in the 5.2 format. But as I mentioned for voice recognition we are using a Windows XP VM. Even if i could get LibreOffice to work on XP, it wouldn't work with Dragon's voice recognition because LO uses Java APIs to draw it's windows, instead of the native Win32 controls. Which have special accessibility methods which Dragon, and other a11y software uses to do it's thing.

Support for LibreOffice exists in a Dragon version 10 service update. But I don't have 10. So this project was about to hit a brick wall. When I remembered that WordPerfect is still shipped to this day. Might there be a Windows XP version of it, and might it work with Dragon I wondered, and might it be able to save in 5.2 format? Turns out that Yes indeed was the answer to all these questions.

In fact i do not know how WordPerfect lost so much market share X3 is better than office 2003 in so many ways it's not even funny.

## The Final Product

![](/assets/images/Screenshot-at-2022-10-01-15-09-02-1024x549.png)

So the final setup involves a 30 year old word processor, running on the latest Fedora, which is also running a VM of a 21 year old operating system , which is in turn running a 16 year old descendant of the 30 year old word processor, and a 14 year old proprietary voice recognition software. The whole setup is in turn replicated to all my machines via my NextCloud Server, And the documents exist on a virtual network drive inside VirtualBox.

To ensure the latest updates to all documents propagate across to my laptop automatically. And Everything is backed up securely.

Some might say this is hilariously over complicated, but I was willing to jump through a lot more hoops to get that 250% productivity boost. In fact you're looking at the smoke test of this system right now.

The TL:DR version of which goes something like. The potential of technology for the disabled will only be realized when, we stop kow towing to an industrial complex which doesn't know us, and is only interested in the money they can get out of us. Also don't give up on your dreams, Hack The system!

Happy almost NaNoWriMo everyone. I
