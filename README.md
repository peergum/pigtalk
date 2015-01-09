pigtalk
=======

<img width="150" bgcolor="#000" src="jesuischarlie.png"/>

A language analysis and interpretation tool

* _Preamble_: This project is NOT RELATED to Hadoop Pig or Pig Latin. "Pig" is just a shortcut for "PeerGum" ;-)

* _Note_: this project is purely experimental at this point. Don't expect anything beautiful or even really usable...

##Program use:
`pigtalk [-d] <filename>` will parse the given file (text) and wait for you to enter the beginning of a word, then suggest the most probable whole word, in a loop.

##Current Progress:

###1. Basic analysis:

Parsing of text to analyse, counting previous and next characters, and determining the spacing character between words.
So far, so good, seems to work fine.

###2. word detection and construction

OK, words are building. I've improved sequencing the characters, building better stats on previous/following chars based on position. It works, but... I need more precision about the whole environment, so the next characters has to depend not just on the previous one, but on all the ones before it in the word. Will patch code accordingly.

Btw, this is not a blog, I'll move it to a proper one asap. Just taking shortcuts for now...

You can type the beginnning of a word ("prefix") and the program will give you his best expectations for the following letters(*). Currently missing a way to give a statistical weight to a word included in another one... Might need to update the prefix table with words and weights...

####-> Current status (2015-01-08):
Updated the prefix/next char stats, so that the program can guess a word from the first character, but not necessarily the longest possible word. E.g: _for_ is more frequent than _form_, so if I type _fo_ I should get the former, not the latter.


###Coming next: 3. sequencing words

Purpose is to detect next word based on previous one(s). We'll use similar stats to character sequences... (coming soon)

_(*) recommended source for stats is language-wikipedia.txt which parses fast and brings useful stats. Despite being much more complete, using text-very-long (Tom Sawyer) is not working well, because the ponctuation is not yet handled properly..._