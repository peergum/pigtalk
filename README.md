pigtalk
=======

A language analysis and interpretation tool

*Preamble*: This project is NOT RELATED to Hadoop Pig or Pig Latin. "Pig" is just a shortcut for "PeerGum" ;-)

##Current Steps:

###1. Basic analysis:

Parsing of text to analyse, counting previous and next characters, and determining the spacing character between words.
So far, so good, seems to work fine.

###2. word detection and construction

OK, words are building. I've improved sequencing the characters, building better stats on previous/following chars based on position. It works, but... I need more precision about the whole environment, so the next characters has to depend not just on the previous one, but on all the ones before it in the word. Will patch code accordingly.

Btw, this is not a blog, I'll move it to a proper one asap. Just taking shortcuts for now...

####-> Current status (2015-01-07):
You can type the beginnning of a word ("prefix") and the program will give you his best expectations for the following letters(*). Currently missing a way to give a statistical weight to a word included in another one... Might need to update the prefix table with words and weights...

###3. sequencing words

Purpose is to detect next word based on previous one(s). We'll use similar stats to character sequences... (coming soon)

_(*) recommended source for stats is language-wikipedia.txt which parses fast and brings useful stats. Despite being much more complete, using text-very-long (Tom Sawyer) is not working well, because the ponctuation is not yet handled properly..._