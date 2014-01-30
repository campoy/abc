/*
This package intends to be a partial parser of the ABC notation as described
in http://abcnotation.com/wiki/abc:standard:v2.1.

The grammar of ABC files I extracted from the specifications linked above is:

tunebook: header? (BLANKLINE (annotation|tune))*
tune: tuneheader tunebody?
tuneheader: refnum \n title (?) \n key

tunebody:

*/
package abc
