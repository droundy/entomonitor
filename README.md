Entomonitor
===========

This is a hack towards writing a distributed bug tracker (DBT).  There
are many distributed bug trackers out there, so you might wonder why I
would write yet another one.  The basic reason is that all the
distributed bug trackers I have seen have one problem or another, and
the nature of distributed bug trackers prevents them from being able
to fix those problems.  In particular, a *good* DBT won't (often?)
change its format on disk, since that would cause problems for its
users.  So you either need to get the format right the first time, or
you need to plan ahead well in the first place.

How to use entomonitor (and status)
-----------------------------------

Below I discuss the design goals for entomonitor, but here I'll give
you the quick intro to how to use it.  But even before the quick
intro:  entomonitor is not yet ready for production use.  It's a fun
experiment, and I'd be delighted to hear that someone is actually
using it, but don't expect a polished or complete product.  I'm not
yet using it myself (except so far as I want to test entomonitor).

You can build entomonitor with gb quite easily, so I won't describe
that here.  If you've never built anything written in go, you probably
want to look elsewhere.  Entomonitor is not a finished tool.

You configure entomonitor by editing the contents of the .entomon
directory.  TODO:  add more documentation here.

Primary goals
-------------

- Bug tracking should take place in the same repository that stores
  source code.  Some bug trackers don't do this
  (e.g. [SD](http://syncwith.us/sd)), but then there is no natural
  synchronization between bug fixes and bug closing.  Putting the data
  with the code means that whatever trickery one does with code (e.g.,
  multiple branches, multiple repositories) will also be naturally
  parallelled in the bug tracker.

- Bug data should be in human-friendly text files.  This is just
  natural for anything that is going to be managed by a distributed
  revision control system.

- Conflicts should never occur.  Users who refrain from editing the
  files manually should never have a conflict show up when merging.
  This is violated, for instance, by
  [ditz](http://ditz.rubyforge.org), which has conflicts if two users
  add comments to the same bug in parallel.  The natural way to avoid
  conflicts is to keep all new data in separate files, and to ensure
  that each new file is (probably) uniquely named.

- There needs to be a command-line interface.

Secondary goals
---------------

- Users should be able to browse bugs via the web.  Static HTML
  generation (as is done by [ditz](http://ditz.rubyforge.org)) would
  be adequate for this, although not optimal.

- You should be able to configure things so users can submit bugs via
  the web.  This requires more trickery, perhaps a stand-alone server
  would be easiest.  But avoiding spam makes things trickier, and
  probably the server should be working on a "scratch" repository that
  is manually synced with the real one.

- It'd be nice to have nice email integration, much like
  [roundup](http://roundup.sourceforge.net).

- At a minimum, email notification would be good.

Various features that are probably desirable or needed
------------------------------------------------------

- Assigning releases or milestones to bugs

- Different classes of bugs (feature requests vs bugs)

- Bug status (closed, open, accepted, etc)

- Bug priorities

- Assigning bugs to specific persons

- Managing a CC list for each bug.

- Allowing users to flexibly configure themselves to be CCed on bugs.
  Perhaps using regular expression matching?

- Manage users and passwords for a web interface?

- Ability to create an "executive summary" for what has changed from
  one version to another

- Nice integration with the version control system, so users can see
  which bugs a present in their version, which ones are present
  in different branches, etc.

Why go?
------

You might wonder why I'm writing this in the go programming language.
I have a couple of reasons.  One is that I enjoy writing in go.
Another is that probably go will be stable well before this tool is,
so its current rate of change shouldn't be an issue.  Also, it's got
some nice features for writing a web server, which seems likely to be
part of this project.  And compiled go code is pretty fast.  That's
probably not an issue, but it'd also be nice to not find that when
I've got 10,000 bugs reported everything begins to slow down.
