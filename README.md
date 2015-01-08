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

###3. sequencing words

Purpose is to detect next word based on previous one(s). We'll use similar stats to character sequences... (coming soon)